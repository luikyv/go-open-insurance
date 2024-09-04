package sec

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/luikyv/go-opf/internal/opinerr"
	"github.com/luikyv/go-opf/internal/resp"
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
	// TODO: Too complicated.
	return func(ctx *gin.Context) {
		tokenInfo := op.TokenInfo(ctx.Writer, ctx.Request)
		if !tokenInfo.IsActive {
			resp.WriteError(ctx, opinerr.New("UNAUTHORISED",
				http.StatusUnauthorized, "invalid token"))
			return
		}

		tokenScopes := strings.Split(tokenInfo.Scopes, " ")
		for _, scope := range scopes {
			matches := false
			for _, tokenScope := range tokenScopes {
				if scope.Matches(tokenScope) {
					matches = true
					break
				}
			}

			if !matches {
				resp.WriteError(ctx, opinerr.New("UNAUTHORISED",
					http.StatusUnauthorized, "invalid token"))
				return
			}
		}

		requestMeta := Meta{
			Subject:  tokenInfo.Subject,
			ClientID: tokenInfo.ClientID,
			Scopes:   tokenScopes,
		}
		exec(ctx, requestMeta)
	}
}
