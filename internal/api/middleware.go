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
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

const (
	headerXFAPIInteractionID = "X-FAPI-Interaction-ID"
	headerIdempotencyID      = "X-Idempotency-Key"
	headerCacheControl       = "Cache-Control"
	headerPragma             = "Pragma"
)

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
				return nil, NewError("ERRO_IDEMPOTENCIA", http.StatusUnprocessableEntity,
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
				return nil, NewError("ERRO_IDEMPOTENCIA", http.StatusUnprocessableEntity,
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
				return nil, NewError(
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

func MetaMiddleware(opinHost string) nethttp.StrictHTTPMiddlewareFunc {
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
			ctx = context.WithValue(ctx, ctxKeyHostURL, opinHost)
			ctx = context.WithValue(ctx, ctxKeyRequestURI, r.URL.RequestURI())
			return f(ctx, w, r, request)
		}
	}
}

// SchemaValidationMiddleware validates incoming requests against an OpenAPI schema.
// If validation fails, it stores the error in the request context.
func SchemaValidationMiddleware(next http.Handler, router routers.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, pathParams, _ := router.FindRoute(r)
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
		}

		ctx := r.Context()
		// TODO: This results in nil pointer dereferencing when the path is not found.
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
	opinErr := NewError("NAO_INFORMADO", http.StatusBadRequest, err.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(opinErr.StatusCode)
	_ = json.NewEncoder(w).Encode(newResponseError(opinErr))
}

func ResponseErrorMiddleware(w http.ResponseWriter, r *http.Request, err error) {
	var opinErr Error
	if !errors.As(err, &opinErr) {
		Logger(r.Context()).Error("unexpected error", slog.Any("error", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrInternal)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(opinErr.StatusCode)
	_ = json.NewEncoder(w).Encode(newResponseError(opinErr))
}

func newResponseError(err Error) ResponseError {
	msg := err.Description
	if len(msg) > 255 {
		msg = msg[:255]
	}
	errData := ErrorInfo{
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
		_ = respErr.Errors.FromErrorInfo(errData)
	} else {
		_ = respErr.Errors.FromErrorInfos([]ErrorInfo{errData})
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
