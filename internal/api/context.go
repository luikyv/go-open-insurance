package api

type ContextKey string

const (
	CtxKeyCorrelationID ContextKey = "correlation_id"
	CtxKeyClientID      ContextKey = "client_id"
	CtxKeySubject       ContextKey = "sub"
)
