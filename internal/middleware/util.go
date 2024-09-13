package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luikyv/go-open-insurance/internal/log"
)

const (
	headerXFAPIInteractionID = "X-FAPI-Interaction-ID"
)

func FAPIID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		interactionID := ctx.GetHeader(headerXFAPIInteractionID)
		if _, err := uuid.Parse(interactionID); err != nil {
			interactionID = uuid.NewString()
		}

		ctx.Set(log.CorrelationIDKey, interactionID)
		ctx.Header(headerXFAPIInteractionID, interactionID)
	}
}

func CacheControl() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Avoid caching.
		ctx.Header("Cache-Control", "no-cache, no-store")
		ctx.Header("Pragma", "no-cache")
	}
}
