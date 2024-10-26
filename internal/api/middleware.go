package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/google/uuid"
	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

const (
	headerXFAPIInteractionID = "X-FAPI-Interaction-ID"
	headerIdempotencyID      = "X-Idempotency-Key"
	headerCacheControl       = "Cache-Control"
	headerPragma             = "Pragma"
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

			token, ok := bearerToken(r)
			if !ok {
				Logger(ctx).Debug("bearer token is required")
				return nil, opinerr.New("UNAUTHORISED", http.StatusUnauthorized,
					"missing token")
			}

			tokenInfo, err := op.TokenInfo(ctx, token)
			if err != nil {
				Logger(ctx).Debug("the token is not active")
				return nil, opinerr.New("UNAUTHORISED", http.StatusUnauthorized,
					"invalid token")
			}

			if err := op.ValidateTokenPoP(
				r,
				token,
				*tokenInfo.Confirmation,
			); err != nil {
				Logger(ctx).Debug("invalid proof of possesion")
				return nil, opinerr.New("UNAUTHORISED", http.StatusUnauthorized,
					"invalid token")
			}

			tokenScopes := strings.Split(tokenInfo.Scopes, " ")
			if !areScopesValid(opts.scopes, tokenScopes) {
				Logger(ctx).Debug("invalid scopes",
					slog.String("token_scopes", tokenInfo.Scopes))
				return nil, opinerr.New("UNAUTHORISED", http.StatusUnauthorized,
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

func bearerToken(r *http.Request) (string, bool) {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return "", false
	}

	tokenParts := strings.Split(tokenHeader, " ")
	if len(tokenParts) != 2 {
		return "", false
	}

	return tokenParts[1], true
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
				return nil, opinerr.New("UNAUTHORIZED", http.StatusUnauthorized,
					"invalid consent")
			}

			return f(ctx, w, r, request)
		}
	}
}

// IdempotencyMiddleware ensures that requests with the same idempotency ID
// are not processed multiple times, returning a cached response if available.
func IdempotencyMiddleware(
	service IdempotencyService,
) nethttp.StrictHTTPMiddlewareFunc {
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
			if !opts.isIdempotent {
				return f(ctx, w, r, request)
			}

			idempotencyID := r.Header.Get(headerIdempotencyID)
			if idempotencyID == "" {
				return nil, opinerr.New("ERRO_IDEMPOTENCIA", http.StatusUnprocessableEntity,
					"missing idempotency id header")
			}

			// Try to fetch a cached response for the idempotency ID.
			idempotentResp, err := service.FetchIdempotencyResponse(
				ctx,
				idempotencyID,
				request,
			)
			// If a cached response exists, write it to the response writer and
			// exit early.
			if err == nil {
				Logger(ctx).Info("return cached idempotency response")
				writeIdempotencyResp(w, r, idempotentResp)
				// returning the response as nil guarantees that the cached
				// response won't be overwritten.
				return nil, nil
			}
			// If the error was not due to "idempotency not found", return an
			// internal error.
			if !errors.Is(err, errIdempotencyNotFound) {
				return nil, opinerr.New("ERRO_IDEMPOTENCIA", http.StatusUnprocessableEntity,
					err.Error())
			}

			// The idempotency record was not found, then process the request
			// and cache the response for next requests with the same idempotency ID.
			response, err = f(ctx, w, r, request)
			if err != nil {
				return nil, err
			}
			_ = service.CreateIdempotency(
				r.Context(),
				idempotencyID,
				request,
				response,
			)
			return response, nil
		}
	}
}

func writeIdempotencyResp(
	w http.ResponseWriter,
	r *http.Request,
	resp string,
) {
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_, _ = w.Write([]byte(resp))
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
			interactionIDIsValid := true
			interactionIDIsRequired := newOperationOptions(operationID).fapiIDIsRequired

			// Verify if the interaction ID is valid, generate a new value if not.
			if _, err := uuid.Parse(interactionID); err != nil {
				interactionIDIsValid = false
				interactionID = uuid.NewString()
			}

			// Return the same interaction ID in the response or a new valid value
			// if the original is invalid.
			w.Header().Add(headerXFAPIInteractionID, interactionID)

			if interactionIDIsRequired && !interactionIDIsValid {
				return nil, opinerr.New(
					"INVALID_INTERACTION_ID",
					http.StatusUnprocessableEntity,
					"The FAPI interaction ID is missing or invalid",
				)
			}

			ctx = context.WithValue(ctx, ctxKeyCorrelationID, interactionID)
			return f(ctx, w, r, request)
		}
	}
}

func MetaMiddleware() nethttp.StrictHTTPMiddlewareFunc {
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
			// TODO: Fill the params here.
			ctx = context.WithValue(ctx, ctxKeyRequestURI, r.URL.RequestURI())
			return f(ctx, w, r, request)
		}
	}
}

func SchemaValidationMiddleware(next http.Handler, router routers.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, pathParams, _ := router.FindRoute(r)
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
		}

		ctx := r.Context()
		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			ctx = context.WithValue(ctx, ctxKeyRequestError, err.Error())
		}
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
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

func RequestErrorMiddleware(w http.ResponseWriter, r *http.Request, err error) {
	Logger(r.Context()).Info("unexpected error", slog.Any("error", err))
	opinErr := opinerr.New("NAO_INFORMADO", http.StatusBadRequest, err.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(opinErr.StatusCode)
	_ = json.NewEncoder(w).Encode(newResponseError(opinErr))
}

func ResponseErrorMiddleware(w http.ResponseWriter, r *http.Request, err error) {
	var opinErr opinerr.Error
	if !errors.As(err, &opinErr) {
		Logger(r.Context()).Error("unexpected error", slog.Any("error", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(opinerr.ErrInternal)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(opinErr.StatusCode)
	_ = json.NewEncoder(w).Encode(newResponseError(opinErr))
}

func newResponseError(err opinerr.Error) ResponseError {
	msg := err.Description
	if len(msg) > 255 {
		msg = msg[:255]
	}
	errData := Error{
		Code:   err.Code,
		Title:  msg,
		Detail: msg,
	}

	respErr := ResponseError{
		Errors: ResponseError_Errors{},
		Meta: &Meta{
			TotalRecords: 1,
			TotalPages:   1,
		},
	}

	if err.StatusCode == http.StatusUnprocessableEntity {
		_ = respErr.Errors.FromError(errData)
	} else {
		_ = respErr.Errors.FromErrors([]Error{errData})
	}

	return respErr
}

// Middleware to disable HTML escaping in responses.
func ResponseEncodingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &responseRecorder{ResponseWriter: w, buf: &bytes.Buffer{}}
		next.ServeHTTP(rec, r)

		modifiedBody := strings.ReplaceAll(rec.buf.String(), `\u0026`, `&`)
		_, _ = w.Write([]byte(modifiedBody))
	})
}

// Custom response recorder to capture the response.
type responseRecorder struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	return rec.buf.Write(b) // Capture response in buffer.
}
