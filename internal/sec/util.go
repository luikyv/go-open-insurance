package sec

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
	"github.com/luikyv/go-open-insurance/internal/opinresp"
)

var (
	errInvalidToken       = opinerr.New("UNAUTHORISED", http.StatusUnauthorized, "invalid token")
	errTokenMissingScopes = opinerr.New("UNAUTHORISED", http.StatusUnauthorized, "token missing scopes")
)

type Meta struct {
	Subject  string
	ClientID string
	Scopes   []string
}

// ProtectedHandler returns a HTTP handler that executes the informed function
// if the request contains the right scopes.
func ProtectedHandler(
	exec func(*gin.Context, Meta),
	op provider.Provider,
	scopes ...goidc.Scope,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenInfo := op.TokenInfo(ctx.Writer, ctx.Request)
		if !tokenInfo.IsActive {
			opinresp.WriteError(ctx, errInvalidToken)
			return
		}

		tokenScopes := strings.Split(tokenInfo.Scopes, " ")
		if !areScopesValid(scopes, tokenScopes) {
			opinresp.WriteError(ctx, errTokenMissingScopes)
			return
		}

		exec(ctx, Meta{
			Subject:  tokenInfo.Subject,
			ClientID: tokenInfo.ClientID,
			Scopes:   tokenScopes,
		})
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
