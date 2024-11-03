package api

import (
	"context"
	"errors"
)

type RequestMeta struct {
	Subject       string
	ClientID      string
	ConsentID     string
	CorrelationID string
	Host          string
	RequestURI    string
	Error         error
}

func (m RequestMeta) RequestURL() string {
	return m.Host + m.RequestURI
}

type ContextKey string

const (
	ctxKeyCorrelationID ContextKey = "correlation_id"
	ctxKeyClientID      ContextKey = "client_id"
	ctxKeySubject       ContextKey = "sub"
	ctxKeyConsentID     ContextKey = "consent_id"
	ctxKeyRequestURI    ContextKey = "request_uri"
	ctxKeyHostURL       ContextKey = "host"
	ctxKeyRequestError  ContextKey = "request_error"
)

func NewRequestMeta(ctx context.Context) RequestMeta {
	meta := RequestMeta{
		ClientID:      strFromCtx(ctx, ctxKeyClientID),
		Subject:       strFromCtx(ctx, ctxKeySubject),
		ConsentID:     strFromCtx(ctx, ctxKeyConsentID),
		CorrelationID: strFromCtx(ctx, ctxKeyCorrelationID),
		RequestURI:    strFromCtx(ctx, ctxKeyRequestURI),
		Host:          strFromCtx(ctx, ctxKeyHostURL),
	}

	if err := strFromCtx(ctx, ctxKeyRequestError); err != "" {
		meta.Error = errors.New(err)
	}

	return meta
}

func strFromCtx(ctx context.Context, key ContextKey) string {
	s := ctx.Value(key)
	if s == nil {
		return ""
	}

	return s.(string)
}
