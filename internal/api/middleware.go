package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/luikyv/go-open-insurance/internal/oidc"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

const (
	headerXFAPIInteractionID = "X-FAPI-Interaction-ID"
	headerCacheControl       = "Cache-Control"
	headerPragma             = "Pragma"
)

var (
	errInvalidToken = opinerr.New("UNAUTHORISED", http.StatusUnauthorized,
		"invalid token")
	errTokenMissingScopes = opinerr.New("UNAUTHORISED", http.StatusUnauthorized,
		"token missing scopes")
	errInvalidConsent = opinerr.New("UNAUTHORISED", http.StatusUnauthorized,
		"invalid consent")
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
			scopes := requiredScopes(operationID)
			if len(scopes) == 0 {
				Logger(ctx).Debug("no scopes are required for the request")
				return f(ctx, w, r, request)
			}

			tokenInfo := op.TokenInfo(w, r)
			if !tokenInfo.IsActive {
				Logger(ctx).Debug("the token is not active")
				return nil, errInvalidToken
			}

			tokenScopes := strings.Split(tokenInfo.Scopes, " ")
			if !areScopesValid(scopes, tokenScopes) {
				Logger(ctx).Debug("invalid scopes",
					slog.String("token_scopes", tokenInfo.Scopes))
				return nil, errTokenMissingScopes
			}

			ctx = context.WithValue(ctx, CtxKeyClientID, tokenInfo.ClientID)
			ctx = context.WithValue(ctx, CtxKeySubject, tokenInfo.Subject)
			consentID, ok := oidc.ConsentID(tokenInfo.Scopes)
			if ok {
				ctx = context.WithValue(ctx, CtxKeyConsentID, consentID)
			}

			return f(ctx, w, r, request)
		}
	}
}

type VerifyPermissionsFunc func(
	ctx context.Context,
	id string,
	permissions ...ConsentPermission,
) error

func AuthPermissionMiddleware(
	verifyPermissions VerifyPermissionsFunc,
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
			permissions := requiredPermissions(operationID)
			if len(permissions) == 0 {
				return f(ctx, w, r, request)
			}

			consentID := ctx.Value(CtxKeyConsentID).(string)
			if err := verifyPermissions(ctx, consentID, permissions...); err != nil {
				Logger(ctx).Debug("the consent is not valid for the request",
					slog.Any("error", err))
				return nil, errInvalidConsent
			}

			return f(ctx, w, r, request)
		}
	}
}

func FAPIIDMiddleware() nethttp.StrictHTTPMiddlewareFunc {
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
			interactionID := r.Header.Get(headerXFAPIInteractionID)
			if _, err := uuid.Parse(interactionID); err != nil {
				interactionID = uuid.NewString()
			}

			w.Header().Add(headerXFAPIInteractionID, interactionID)
			ctx = context.WithValue(ctx, CtxKeyCorrelationID, interactionID)
			return f(ctx, w, r, request)
		}
	}
}

func CacheControlMiddleware() nethttp.StrictHTTPMiddlewareFunc {
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
			w.Header().Add(headerCacheControl, "no-cache, no-store")
			w.Header().Add(headerPragma, "no-cache")
			return f(ctx, w, r, request)
		}
	}
}

func ValidationErrorHandler() nethttpmiddleware.ErrorHandler {
	return func(w http.ResponseWriter, message string, statusCode int) {
		opinErr := opinerr.New("INVALID_REQUEST", http.StatusBadRequest, message)
		w.WriteHeader(opinErr.StatusCode)
		_ = json.NewEncoder(w).Encode(newResponseError(opinErr))
	}
}

func ResponseErrorMiddleware(w http.ResponseWriter, r *http.Request, err error) {
	var opinErr opinerr.Error
	if !errors.As(err, &opinErr) {
		Logger(r.Context()).Error("unexpected error", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(opinerr.ErrInternal)
	}

	w.WriteHeader(opinErr.StatusCode)
	_ = json.NewEncoder(w).Encode(newResponseError(opinErr))
}

func newResponseError(err opinerr.Error) ResponseError {
	title := err.Description
	if len(title) > 255 {
		title = title[:255]
	}

	detail := err.Description
	if len(detail) > 2048 {
		detail = detail[:2048]
	}
	return ResponseError{
		Errors: []Error{
			{
				Code:   err.Code,
				Title:  title,
				Detail: detail,
			},
		},
		Meta: &Meta{
			TotalRecords: 1,
			TotalPages:   1,
		},
	}
}
