package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/go-jose/go-jose/v4"
	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/consent"
	"github.com/luikyv/go-open-insurance/internal/oidc"
	"github.com/luikyv/go-open-insurance/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	headerClientCert = "X-Client-Cert"
)

func openidProvider(
	db *mongo.Database,
	userService user.Service,
	consentService consent.Service,
	host, mtlsHost, prefix string,
) (
	provider.Provider,
	error,
) {

	// Get the file path of the source file.
	_, filename, _, _ := runtime.Caller(0)
	sourceDir := filepath.Dir(filename)

	keysDir := filepath.Join(sourceDir, "../../keys")
	templatesDirPath := filepath.Join(sourceDir, "../../templates")

	return provider.New(
		goidc.ProfileOpenID,
		host,
		privateJWKS(filepath.Join(keysDir, "server.jwks")),
		provider.WithPathPrefix(prefix),
		provider.WithClientStorage(oidc.NewClientManager(db)),
		provider.WithAuthnSessionStorage(oidc.NewAuthnSessionManager(db)),
		provider.WithGrantSessionStorage(oidc.NewGrantSessionManager(db)),
		provider.WithTokenAuthnMethods(goidc.ClientAuthnPrivateKeyJWT),
		provider.WithScopes(api.Scopes...),
		provider.WithMTLS(mtlsHost, clientCertFunc()),
		provider.WithTLSCertTokenBindingRequired(),
		provider.WithPAR(60),
		provider.WithJAR(jose.PS256),
		provider.WithJAREncryption(jose.RSA_OAEP),
		provider.WithJARContentEncryptionAlgs(jose.A256GCM),
		provider.WithJARM(jose.PS256),
		provider.WithPrivateKeyJWTSignatureAlgs(jose.PS256),
		provider.WithIssuerResponseParameter(),
		provider.WithClaimsParameter(),
		provider.WithClaims(goidc.ClaimEmail, goidc.ClaimEmailVerified),
		provider.WithPKCE(goidc.CodeChallengeMethodSHA256),
		provider.WithAuthorizationCodeGrant(),
		provider.WithImplicitGrant(),
		provider.WithRefreshTokenGrant(shoudIssueRefreshTokenFunc(), 600),
		provider.WithACRs(api.ACROpenInsuranceLOA2, api.ACROpenInsuranceLOA3),
		provider.WithUserInfoEncryption(jose.RSA_OAEP),
		provider.WithStaticClient(client("client_one", keysDir)),
		provider.WithStaticClient(client("client_two", keysDir)),
		provider.WithHandleGrantFunc(handleGrantFunc(consentService)),
		provider.WithPolicy(oidc.Policy(templatesDirPath, host+prefix, userService, consentService)),
		provider.WithNotifyErrorFunc(logErrorFun()),
	)
}

func handleGrantFunc(consentService consent.Service) goidc.HandleGrantFunc {
	return func(r *http.Request, gi *goidc.GrantInfo) error {
		consentID, ok := api.ConsentID(gi.ActiveScopes)
		if !ok {
			return nil
		}

		meta := api.RequestMeta{
			ClientID: gi.ClientID,
		}
		consent, err := consentService.Get(r.Context(), meta, consentID)
		if err != nil {
			return err
		}

		if !consent.IsAuthorized() {
			return goidc.NewError(goidc.ErrorCodeInvalidRequest,
				"consent is not authorized")
		}

		return nil
	}
}

func shoudIssueRefreshTokenFunc() goidc.ShouldIssueRefreshTokenFunc {
	return func(client *goidc.Client, grantInfo goidc.GrantInfo) bool {
		return slices.Contains(client.GrantTypes, goidc.GrantRefreshToken)
	}
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

func clientCertFunc() goidc.ClientCertFunc {
	return func(r *http.Request) (*x509.Certificate, error) {
		rawClientCert := r.Header.Get(headerClientCert)
		if rawClientCert == "" {
			return nil, errors.New("the client certificate was not informed")
		}

		// Apply URL decoding.
		rawClientCert, err := url.QueryUnescape(rawClientCert)
		if err != nil {
			return nil, fmt.Errorf("could not url decode the client certificate: %w", err)
		}

		clientCertPEM, _ := pem.Decode([]byte(rawClientCert))
		if clientCertPEM == nil {
			return nil, errors.New("could not decode the client certificate")
		}

		clientCert, err := x509.ParseCertificate(clientCertPEM.Bytes)
		if err != nil {
			return nil, fmt.Errorf("could not parse the client certificate: %w", err)
		}

		return clientCert, nil
	}
}

func logErrorFun() goidc.NotifyErrorFunc {
	return func(r *http.Request, err error) {
		api.Logger(r.Context()).Info(
			"error during request",
			slog.String("uri", r.URL.RequestURI()),
			slog.String("error", err.Error()),
		)
	}
}
