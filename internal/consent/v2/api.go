package consentv2

import (
	"context"
	"fmt"
	"time"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/consent"
)

type Server struct {
	baseURL string
	service consent.Service
}

func NewServer(baseURL string, service consent.Service) Server {
	return Server{
		baseURL: baseURL,
		service: service,
	}
}

func (s Server) CreateConsentV2(
	ctx context.Context,
	request api.CreateConsentV2RequestObject,
) (
	api.CreateConsentV2ResponseObject,
	error,
) {

	consent := newConsent(ctx, *request.Body)
	if err := s.service.Create(ctx, consent); err != nil {
		return nil, err
	}

	resp := newResponse(consent, s.baseURL)
	return api.CreateConsentV2201JSONResponse(resp), nil
}

func (s Server) DeleteConsentV2(
	ctx context.Context,
	request api.DeleteConsentV2RequestObject,
) (
	api.DeleteConsentV2ResponseObject,
	error,
) {

	if err := s.service.Reject(
		ctx,
		request.ConsentId,
		consent.RejectionInfo{
			RejectedBy: api.ConsentRejectedByTPP,
			Reason:     api.ConsentRejectedReasonCodeCUSTOMERMANUALLYREVOKED,
		}); err != nil {
		return nil, err
	}

	return api.DeleteConsentV2204Response{}, nil
}

func (s Server) ConsentV2(
	ctx context.Context,
	request api.ConsentV2RequestObject,
) (
	api.ConsentV2ResponseObject,
	error,
) {
	consent, err := s.service.Get(ctx, request.ConsentId)
	if err != nil {
		return nil, err
	}

	resp := newResponse(consent, s.baseURL)
	return api.ConsentV2200JSONResponse(resp), nil
}

func newConsent(ctx context.Context, req api.CreateConsentRequestV2) consent.Consent {
	clientID := ctx.Value(api.CtxKeyClientID).(string)
	now := time.Now().UTC()
	consent := consent.Consent{
		ID:          consent.ID(),
		Status:      api.ConsentStatusAWAITINGAUTHORISATION,
		UserCPF:     req.Data.LoggedUser.Document.Identification,
		ClientId:    clientID,
		Permissions: req.Data.Permissions,
		CreatedAt:   now,
		UpdatedAt:   now,
		ExpiresAt:   req.Data.ExpirationDateTime.Time,
	}

	if req.Data.BusinessEntity != nil {
		consent.BusinessCNPJ = req.Data.BusinessEntity.Document.Identification
	}

	return consent
}

func newResponse(consent consent.Consent, baseURL string) api.ConsentResponseV2 {

	data := api.ConsentResponseDataV2{
		ConsentId:            consent.ID,
		Status:               consent.Status,
		Permissions:          consent.Permissions,
		CreationDateTime:     api.NewDateTime(consent.CreatedAt),
		StatusUpdateDateTime: api.NewDateTime(consent.CreatedAt),
		ExpirationDateTime:   api.NewDateTime(consent.ExpiresAt),
	}

	if consent.RejectionInfo != nil {
		data.Rejection = &api.ConsentRejection{
			RejectedBy: consent.RejectionInfo.RejectedBy,
			Reason: api.ConsentRejectedReason{
				Code: consent.RejectionInfo.Reason,
			},
		}
	}

	return api.ConsentResponseV2{
		Data: data,
		Meta: &api.Meta{
			TotalRecords: 1,
			TotalPages:   1,
		},
		Links: &api.Links{
			Self: fmt.Sprintf("%s/consents/v2/consents/%s", baseURL, consent.ID),
		},
	}
}