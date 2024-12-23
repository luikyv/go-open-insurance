package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

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
			opts := newOperationOptions(operationID)
			if len(opts.scopes) == 0 {
				Logger(ctx).Debug("no scopes are required for the request")
				return f(ctx, w, r, request)
			}

			tokenInfo, err := op.TokenInfoFromRequest(w, r)
			if err != nil {
				Logger(ctx).Debug("the token is not active")
				return nil, NewError("UNAUTHORISED", http.StatusUnauthorized,
					"invalid token")
			}

			tokenScopes := strings.Split(tokenInfo.Scopes, " ")
			if !areScopesValid(opts.scopes, tokenScopes) {
				Logger(ctx).Debug("invalid scopes",
					slog.String("token_scopes", tokenInfo.Scopes))
				return nil, NewError("UNAUTHORISED", http.StatusUnauthorized,
					"token missing scopes")
			}

			ctx = context.WithValue(ctx, ctxKeyClientID, tokenInfo.ClientID)
			ctx = context.WithValue(ctx, ctxKeySubject, tokenInfo.Subject)
			consentID, ok := ConsentID(tokenInfo.Scopes)
			if ok {
				ctx = context.WithValue(ctx, ctxKeyConsentID, consentID)
			}

			return f(ctx, w, r, request)
		}
	}
}

func AuthPermissionMiddleware(
	consentService interface {
		Verify(
			ctx context.Context,
			meta RequestMeta,
			consentID string,
			permissions ...ConsentPermission,
		) error
	},
) StrictMiddlewareFunc {
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
			opts := newOperationOptions(operationID)
			if len(opts.permissions) == 0 {
				return f(ctx, w, r, request)
			}

			meta := NewRequestMeta(ctx)
			if err = consentService.Verify(ctx, meta, meta.ConsentID, opts.permissions...); err != nil {
				Logger(ctx).Debug("the consent is not valid for the request",
					slog.Any("error", err))
				return nil, NewError("UNAUTHORIZED", http.StatusForbidden,
					"invalid consent")
			}

			return f(ctx, w, r, request)
		}
	}
}
