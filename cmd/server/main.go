package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/go-jose/go-jose/v4"
	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/capitalizationtitle"
	"github.com/luikyv/go-open-insurance/internal/consent"
	"github.com/luikyv/go-open-insurance/internal/customer"
	"github.com/luikyv/go-open-insurance/internal/endorsement"
	"github.com/luikyv/go-open-insurance/internal/oidc"
	"github.com/luikyv/go-open-insurance/internal/quoteauto"
	"github.com/luikyv/go-open-insurance/internal/resource"
	"github.com/luikyv/go-open-insurance/internal/user"
	"github.com/luikyv/go-open-insurance/internal/webhook"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	dbSchema              = getEnv("MOCKIN_DB_SCHEMA", "mockin")
	dbStringConnection    = getEnv("MOCKIN_DB_CONNECTION", "mongodb://localhost:27017/mockin")
	port                  = getEnv("MOCKIN_PORT", "80")
	awsBaseEndpoint       = getEnv("MOCKIN_AWS_BASE_ENDPOINT", "http://localhost:4566")
	host                  = getEnv("MOCKIN_HOST", "https://mockin.local")
	mtlsHost              = getEnv("MOCKIN_MTLS_HOST", "https://matls-mockin.local")
	kmsSigningKeyAlias    = getEnv("MOCKIN_KMS_SIGNING_KEY_ALIAS", "alias/mockin/signing-key")
	kmsEncryptionKeyAlias = getEnv("MOCKIN_KMS_ENCRYPTION_KEY_ALIAS", "alias/mockin/encryption-key")
	apiPrefixOIDC         = "/auth"
	apiPrefixOPIN         = "/open-insurance"
)

type ConsentServerV2 = consent.ServerV2
type CustomerServerV1 = customer.ServerV1
type ResourceServerV2 = resource.ServerV2
type CapitalizationTitleServerV1 = capitalizationtitle.ServerV1
type EndorsementServerV1 = endorsement.ServerV1
type QuoteAutoServerV1 = quoteauto.ServerV1
type opinServer struct {
	ConsentServerV2
	CustomerServerV1
	ResourceServerV2
	CapitalizationTitleServerV1
	EndorsementServerV1
	QuoteAutoServerV1
}

func main() {
	kmsClient := kmsClient()
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}

	// Storage.
	userStorage := user.NewStorage()
	consentStorage := consent.NewStorage(db)
	idempotencyStorage := api.NewIdempotencyStorage(db)
	resourceStorage := resource.NewStorage()
	customerStorage := customer.NewStorage()
	capTitleStorage := capitalizationtitle.NewStorage()
	quoteAutoStorage := quoteauto.NewStorage(db)

	// Services.
	userService := user.NewService(userStorage)
	consentService := consent.NewService(consentStorage, userService)
	// OpenID Provider.
	op, err := openidProvider(db, kmsClient, userService, consentService)
	if err != nil {
		log.Fatal(err)
	}
	webhookService := webhook.NewService(op, httpClientFunc())
	idempotencyService := api.NewIdempotencyService(idempotencyStorage)
	resourceService := resource.NewService(resourceStorage, consentService)
	customerService := customer.NewService(customerStorage)
	capitalizationtitleService := capitalizationtitle.NewService(capTitleStorage, resourceService)
	endorsementService := endorsement.NewService(consentService, resourceService)
	quoteAutoService := quoteauto.NewService(quoteAutoStorage, webhookService)

	// Server.
	server := opinServer{
		ConsentServerV2:             consent.NewServerV2(consentService),
		CustomerServerV1:            customer.NewServerV1(customerService),
		ResourceServerV2:            resource.NewServerV2(resourceService),
		CapitalizationTitleServerV1: capitalizationtitle.NewServerV1(capitalizationtitleService),
		EndorsementServerV1:         endorsement.NewServerV1(endorsementService),
		QuoteAutoServerV1:           quoteauto.NewServerV1(quoteAutoService),
	}

	strictHandler := api.NewStrictHandlerWithOptions(
		server,
		[]nethttp.StrictHTTPMiddlewareFunc{
			api.MetaMiddleware(mtlsHost),
			api.CacheControlMiddleware(),
			api.AuthPermissionMiddleware(consentService),
			api.AuthScopeMiddleware(op),
			api.IdempotencyMiddleware(idempotencyService),
			api.FAPIIDMiddleware(),
		},
		api.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  api.RequestErrorMiddleware,
			ResponseErrorHandlerFunc: api.ResponseErrorMiddleware,
		},
	)

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}
	router, err := gorillamux.NewRouter(swagger)
	if err != nil {
		log.Fatal(err)
	}

	opinHandler := api.HandlerFromMux(strictHandler, http.NewServeMux())
	opinHandler = api.SchemaValidationMiddleware(opinHandler, router)
	opinHandler = api.ResponseEncodingMiddleware(opinHandler)

	mux := http.NewServeMux()
	mux.Handle(apiPrefixOIDC+"/", op.Handler())
	mux.Handle(apiPrefixOPIN+"/", opinHandler)

	// Run.
	if err := loadMocks(
		userService,
		customerService,
		resourceService,
		capitalizationtitleService,
	); err != nil {
		log.Fatal(err)
	}
	s := &http.Server{
		Handler: mux,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func dbConnection() (*mongo.Database, error) {
	ctx := context.Background()

	conn, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(dbStringConnection).SetBSONOptions(&options.BSONOptions{
			UseJSONStructTags: true,
			NilMapAsEmpty:     true,
			NilSliceAsEmpty:   true,
		}),
	)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return conn.Database(dbSchema), nil
}

func openidProvider(
	db *mongo.Database,
	kmsClient *kms.Client,
	userService user.Service,
	consentService consent.Service,
) (
	provider.Provider,
	error,
) {

	// Get the file path of the source file.
	_, filename, _, _ := runtime.Caller(0)
	sourceDir := filepath.Dir(filename)

	// TODO: This will cause problems for the docker file.
	keysDir := filepath.Join(sourceDir, "../../keys")
	templatesDirPath := filepath.Join(sourceDir, "../../templates")

	return provider.New(
		goidc.ProfileOpenID,
		host,
		oidc.JWKSFunc(kmsClient, kmsSigningKeyAlias, kmsEncryptionKeyAlias),
		provider.WithSignFunc(oidc.SignFunc(kmsClient, kmsSigningKeyAlias)),
		provider.WithDecryptFunc(oidc.DecryptFunc(kmsClient, kmsEncryptionKeyAlias)),
		provider.WithPathPrefix(apiPrefixOIDC),
		provider.WithClientStorage(oidc.NewClientManager(db)),
		provider.WithAuthnSessionStorage(oidc.NewAuthnSessionManager(db)),
		provider.WithGrantSessionStorage(oidc.NewGrantSessionManager(db)),
		provider.WithScopes(api.Scopes...),
		provider.WithTokenOptions(oidc.TokenOptionsFunc()),
		provider.WithAuthorizationCodeGrant(),
		provider.WithImplicitGrant(),
		provider.WithRefreshTokenGrant(oidc.ShoudIssueRefreshTokenFunc(), 600),
		provider.WithClientCredentialsGrant(),
		provider.WithTokenAuthnMethods(goidc.ClientAuthnPrivateKeyJWT),
		provider.WithPrivateKeyJWTSignatureAlgs(jose.PS256),
		provider.WithMTLS(mtlsHost, oidc.ClientCertFunc()),
		provider.WithTLSCertTokenBindingRequired(),
		provider.WithPAR(60),
		provider.WithJAR(jose.PS256),
		provider.WithJAREncryption(jose.RSA_OAEP),
		provider.WithJARContentEncryptionAlgs(jose.A256GCM),
		provider.WithJARM(jose.PS256),
		provider.WithIssuerResponseParameter(),
		provider.WithPKCE(goidc.CodeChallengeMethodSHA256),
		provider.WithACRs(api.ACROpenInsuranceLOA2, api.ACROpenInsuranceLOA3),
		provider.WithUserSignatureAlgs(jose.PS256),
		provider.WithUserInfoEncryption(jose.RSA_OAEP),
		provider.WithStaticClient(client("client_one", keysDir)),
		provider.WithStaticClient(client("client_two", keysDir)),
		provider.WithHandleGrantFunc(oidc.HandleGrantFunc(consentService)),
		provider.WithPolicy(oidc.Policy(templatesDirPath, host+apiPrefixOIDC, userService, consentService)),
		provider.WithNotifyErrorFunc(oidc.LogErrorFunc()),
		provider.WithDCR(
			oidc.DCRFunc(api.Scopes),
			func(r *http.Request, s string) error {
				return nil
			},
		),
		provider.WithHTTPClientFunc(httpClientFunc()),
	)
}

func client(clientID string, keysDir string) *goidc.Client {
	var scopes []string
	for _, scope := range api.Scopes {
		scopes = append(scopes, scope.ID)
	}

	privateJWKS := privateJWKS(filepath.Join(keysDir, clientID+".jwks"))
	publicJWKS := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
	for _, jwk := range privateJWKS.Keys {
		publicJWKS.Keys = append(publicJWKS.Keys, jwk.Public())
	}
	rawPublicJWKS, _ := json.Marshal(publicJWKS)
	return &goidc.Client{
		ID: clientID,
		ClientMetaInfo: goidc.ClientMetaInfo{
			TokenAuthnMethod: goidc.ClientAuthnPrivateKeyJWT,
			ScopeIDs:         strings.Join(scopes, " "),
			RedirectURIs: []string{
				"https://localhost.emobix.co.uk:8443/test/a/mockin/callback",
			},
			GrantTypes: []goidc.GrantType{
				goidc.GrantAuthorizationCode,
				goidc.GrantRefreshToken,
				goidc.GrantClientCredentials,
				goidc.GrantImplicit,
			},
			ResponseTypes: []goidc.ResponseType{
				goidc.ResponseTypeCode,
				goidc.ResponseTypeCodeAndIDToken,
			},
			PublicJWKS:           rawPublicJWKS,
			IDTokenKeyEncAlg:     jose.RSA_OAEP,
			IDTokenContentEncAlg: jose.A128CBC_HS256,
		},
	}
}

func privateJWKS(filePath string) jose.JSONWebKeySet {
	absPath, _ := filepath.Abs(filePath)
	jwksFile, err := os.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer jwksFile.Close()

	jwksBytes, err := io.ReadAll(jwksFile)
	if err != nil {
		log.Fatal(err)
	}

	var jwks jose.JSONWebKeySet
	if err := json.Unmarshal(jwksBytes, &jwks); err != nil {
		log.Fatal(err)
	}

	return jwks
}

func httpClientFunc() goidc.HTTPClientFunc {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Renegotiation:      tls.RenegotiateOnceAsClient,
				InsecureSkipVerify: true,
			},
		},
	}

	return func(ctx context.Context) *http.Client {
		return client
	}
}

func kmsClient() *kms.Client {
	return kms.New(kms.Options{
		BaseEndpoint: &awsBaseEndpoint,
	})
}

// getEnv retrieves an environment variable or returns a fallback value if not found
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
