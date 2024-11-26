package oidc

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/go-jose/go-jose/v4"
	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/consent"
)

const (
	HeaderClientCert = "X-Client-Cert"
)

func HandleGrantFunc(consentService consent.Service) goidc.HandleGrantFunc {
	return func(r *http.Request, gi *goidc.GrantInfo) error {
		consentID, ok := api.ConsentID(gi.ActiveScopes)
		if !ok {
			return nil
		}

		meta := api.RequestMeta{
			ClientID: gi.ClientID,
		}
		consent, err := consentService.Fetch(r.Context(), meta, consentID)
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

func ShoudIssueRefreshTokenFunc() goidc.ShouldIssueRefreshTokenFunc {
	return func(client *goidc.Client, grantInfo goidc.GrantInfo) bool {
		return slices.Contains(client.GrantTypes, goidc.GrantRefreshToken) &&
			(grantInfo.GrantType == goidc.GrantAuthorizationCode || grantInfo.GrantType == goidc.GrantRefreshToken)
	}
}

func ClientCertFunc() goidc.ClientCertFunc {
	return func(r *http.Request) (*x509.Certificate, error) {
		rawClientCert := r.Header.Get(HeaderClientCert)
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

func LogErrorFunc() goidc.NotifyErrorFunc {
	return func(ctx context.Context, err error) {
		api.Logger(ctx).Info(
			"error during request",
			slog.String("error", err.Error()),
		)
	}
}

func DCRFunc(scopes []goidc.Scope) goidc.HandleDynamicClientFunc {
	var scopeIDs []string
	for _, scope := range scopes {
		scopeIDs = append(scopeIDs, scope.ID)
	}
	scopeIDsStr := strings.Join(scopeIDs, " ")
	return func(r *http.Request, c *goidc.ClientMetaInfo) error {
		c.ScopeIDs = scopeIDsStr
		return nil
	}
}

func TokenOptionsFunc() goidc.TokenOptionsFunc {
	return func(gi goidc.GrantInfo, c *goidc.Client) goidc.TokenOptions {
		return goidc.NewJWTTokenOptions(jose.PS256, 300)
	}
}
