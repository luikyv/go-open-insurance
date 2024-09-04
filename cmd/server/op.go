package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-jose/go-jose/v4"
	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/luikyv/go-opf/internal/authn"
	"github.com/luikyv/go-opf/internal/consent"
	"github.com/luikyv/go-opf/internal/oidc"
	"github.com/luikyv/go-opf/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
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
	authenticator := authn.New(userService, consentService, host+prefix)
	ps256ServerKeyID := "ps256_key"
	return provider.New(
		host,
		privateJWKS("../../keys/server_jwks.json"),
		provider.WithPathPrefix(prefix),
		provider.WithClientStorage(oidc.NewClientManager(db)),
		provider.WithAuthnSessionStorage(oidc.NewAuthnSessionManager(db)),
		provider.WithGrantSessionStorage(oidc.NewGrantSessionManager(db)),
		provider.WithScopes(oidc.Scopes...),
		provider.WithMTLS(mtlsHost),
		provider.WithJAR(jose.PS256),
		provider.WithJAREncryption("enc_key"),
		provider.WithJARContentEncryptionAlgs(jose.A256GCM),
		provider.WithJARM(ps256ServerKeyID),
		provider.WithPrivateKeyJWTAuthn(jose.PS256),
		provider.WithIssuerResponseParameter(),
		provider.WithClaimsParameter(),
		provider.WithClaims(goidc.ClaimEmail, goidc.ClaimEmailVerified),
		provider.WithDPoP(jose.PS256, jose.ES256),
		provider.WithPKCE(goidc.CodeChallengeMethodSHA256),
		provider.WithRefreshTokenGrant(),
		provider.WithACRs(oidc.ACROpenInsuranceLOA2, oidc.ACROpenInsuranceLOA3),
		provider.WithDCR(dcrFunc()),
		provider.WithTokenOptions(tokenOptionFunc(ps256ServerKeyID)),
		provider.WithUserInfoEncryption(jose.RSA_OAEP_256),
		provider.WithStaticClient(client("client_one")),
		provider.WithStaticClient(client("client_two")),
		provider.WithPolicy(goidc.NewPolicy(
			"policy",
			func(r *http.Request, c *goidc.Client, as *goidc.AuthnSession) bool {
				return true
			},
			authenticator.Authenticate,
		)),
	)
}

func dcrFunc() goidc.HandleDynamicClientFunc {
	var scopes []string
	for _, scope := range oidc.Scopes {
		scopes = append(scopes, scope.ID)
	}
	scopeStr := strings.Join(scopes, " ")
	return func(r *http.Request, c *goidc.ClientMetaInfo) error {
		c.ScopeIDs = scopeStr
		return nil
	}
}

func tokenOptionFunc(keyID string) goidc.TokenOptionsFunc {
	return func(client *goidc.Client, scopes string) (goidc.TokenOptions, error) {
		return goidc.NewJWTTokenOptions(keyID, 600), nil
	}
}

func client(clientID string) *goidc.Client {

	var scopes []string
	for _, scope := range oidc.Scopes {
		scopes = append(scopes, scope.ID)
	}

	privateJWKS := privateJWKS(fmt.Sprintf("../../keys/%s_jwks.json", clientID))
	publicJWKS := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
	for _, jwk := range privateJWKS.Keys {
		publicJWKS.Keys = append(publicJWKS.Keys, jwk.Public())
	}
	rawPublicJWKS, _ := json.Marshal(publicJWKS)
	return &goidc.Client{
		ID: clientID,
		ClientMetaInfo: goidc.ClientMetaInfo{
			AuthnMethod: goidc.ClientAuthnPrivateKeyJWT,
			ScopeIDs:    strings.Join(scopes, " "),
			RedirectURIs: []string{
				"https://localhost.emobix.co.uk:8443/test/a/gopin/callback",
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
		panic(err.Error())
	}
	defer jwksFile.Close()

	jwksBytes, err := io.ReadAll(jwksFile)
	if err != nil {
		panic(err.Error())
	}

	var jwks jose.JSONWebKeySet
	if err := json.Unmarshal(jwksBytes, &jwks); err != nil {
		panic(err.Error())
	}

	return jwks
}
