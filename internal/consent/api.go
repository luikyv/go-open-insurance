package consent

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/luikyv/go-open-insurance/internal/oidc"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
	"github.com/luikyv/go-open-insurance/internal/opinresp"
	"github.com/luikyv/go-open-insurance/internal/sec"
)

type Router struct {
	baseURL   string
	nameSpace string
	provider  provider.Provider
	service   Service
}

func NewRouter(
	op provider.Provider,
	service Service,
	baseURL, nameSpace string,
) Router {
	return Router{
		baseURL:   baseURL,
		nameSpace: nameSpace,
		provider:  op,
		service:   service,
	}
}

func (r Router) AddRoutesV2(
	router gin.IRouter,
) {
	consentRouter := router.Group(apiPrefixConsentsV2)

	consentRouter.POST("/consents",
		sec.ProtectedHandler(r.handlePostV2, r.provider, oidc.ScopeConsents))

	consentRouter.GET("/consents/:consent_id",
		sec.ProtectedHandler(r.handleGetV2, r.provider, oidc.ScopeConsents))

	consentRouter.DELETE("/consents/:consent_id",
		sec.ProtectedHandler(r.handleDeleteV2, r.provider, oidc.ScopeConsents))
}

func (r Router) handlePostV2(ctx *gin.Context, meta sec.Meta) {
	var req struct {
		Data requestData `json:"data"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		opinresp.WriteError(ctx, opinerr.New("INVALID_REQUEST", http.StatusBadRequest, err.Error()))
		return
	}

	consent := newConsent(req.Data, meta.ClientID, r.nameSpace)
	if err := r.service.Create(ctx, meta, consent); err != nil {
		opinresp.WriteError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, consent.response(r.baseURL+apiPrefixConsentsV2))
}

func (r Router) handleGetV2(ctx *gin.Context, meta sec.Meta) {
	consentID := ctx.Param("consent_id")
	consent, err := r.service.Get(ctx, meta, consentID)
	if err != nil {
		opinresp.WriteError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, consent.response(r.baseURL+apiPrefixConsentsV2))
}

func (r Router) handleDeleteV2(ctx *gin.Context, meta sec.Meta) {
	consentID := ctx.Param("consent_id")
	if err := r.service.Reject(
		ctx,
		meta,
		consentID,
		RejectionInfo{
			RejectedBy: RejectedByUser,
			Reason:     RejectionReasonCustomerManuallyRejected,
		},
	); err != nil {
		opinresp.WriteError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
