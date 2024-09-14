package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/luikyv/go-open-insurance/internal/oidc"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

const (
	headerXFAPIInteractionID = "X-FAPI-Interaction-ID"
	headerCacheControl       = "Cache-Control"
	headerPragma             = "Pragma"
)

var (
	errInvalidToken       = opinerr.New("UNAUTHORISED", http.StatusUnauthorized, "invalid token")
	errTokenMissingScopes = opinerr.New("UNAUTHORISED", http.StatusUnauthorized, "token missing scopes")
)

func FAPIIDMiddleware(
	f nethttp.StrictHTTPHandlerFunc,
	operationID string,
) nethttp.StrictHTTPHandlerFunc {
	return func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		request interface{},
	) (
		response interface{},
		err error,
	) {
		interactionID := r.Header.Get(headerXFAPIInteractionID)
		if _, err := uuid.Parse(interactionID); err != nil {
			interactionID = uuid.NewString()
		}

		w.Header().Add(headerXFAPIInteractionID, interactionID)
		ctx = context.WithValue(ctx, CtxKeyCorrelationID, interactionID)
		return f(ctx, w, r, request)
	}
}

func CacheControlMiddleware(
	f nethttp.StrictHTTPHandlerFunc,
	operationID string,
) nethttp.StrictHTTPHandlerFunc {
	return func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		request interface{},
	) (
		response interface{},
		err error,
	) {

		w.Header().Add(headerCacheControl, "no-cache, no-store")
		w.Header().Add(headerPragma, "no-cache")
		return f(ctx, w, r, request)
	}
}

func AuthScopeMiddleware(op provider.Provider) StrictMiddlewareFunc {
	return func(
		f nethttp.StrictHTTPHandlerFunc,
		operationID string,
	) nethttp.StrictHTTPHandlerFunc {
		return func(
			ctx context.Context,
			w http.ResponseWriter,
			r *http.Request,
			request interface{},
		) (
			response interface{},
			err error,
		) {
			scopes := requiredScopes(operationID)
			if len(scopes) == 0 {
				return f(ctx, w, r, request)
			}

			tokenInfo := op.TokenInfo(w, r)
			if !tokenInfo.IsActive {
				return nil, errInvalidToken
			}

			tokenScopes := strings.Split(tokenInfo.Scopes, " ")
			if !areScopesValid(scopes, tokenScopes) {
				return nil, errTokenMissingScopes
			}

			ctx = context.WithValue(ctx, CtxKeyClientID, tokenInfo.ClientID)
			ctx = context.WithValue(ctx, CtxKeySubject, tokenInfo.Subject)
			return f(ctx, w, r, request)
		}
	}
}

func ErrorMiddleware(w http.ResponseWriter, r *http.Request, err error) {

}

func requiredScopes(operationID string) []goidc.Scope {
	switch operationID {
	case "CreateConsentV2":
		return []goidc.Scope{oidc.ScopeConsents}
	case "ConsentV2":
		return []goidc.Scope{oidc.ScopeConsents}
	case "DeleteConsentV2":
		return []goidc.Scope{oidc.ScopeConsents}
	default:
		return nil
	}
}

// areScopesValid verifies every scope in requiredScopes has a match among
// scopes.
// scopes can have more scopes than the defined at requiredScopes, but the
// contrary results in false.
func areScopesValid(requiredScopes []goidc.Scope, scopes []string) bool {
	for _, requiredScope := range requiredScopes {
		if !isScopeValid(requiredScope, scopes) {
			return false
		}
	}
	return true
}

// isScopeValid verifies if requireScope has a match in scopes.
func isScopeValid(requiredScope goidc.Scope, scopes []string) bool {
	for _, scope := range scopes {
		if requiredScope.Matches(scope) {
			return true
		}
	}

	return false
}
