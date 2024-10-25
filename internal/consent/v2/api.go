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

func NewServer(
	baseURL string,
	service consent.Service,
) Server {
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

	meta := api.NewRequestMeta(ctx)
	consent := newConsent(meta, *request.Body)
	if err := s.service.Create(ctx, meta, consent); err != nil {
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
	meta := api.NewRequestMeta(ctx)
	c, err := s.service.Get(ctx, meta, request.ConsentId)
	if err != nil {
		return nil, err
	}

	reason := api.ConsentRejectedReasonCodeCUSTOMERMANUALLYREJECTED
	if c.IsAuthorized() {
		reason = api.ConsentRejectedReasonCodeCUSTOMERMANUALLYREVOKED
	}

	if err := s.service.Reject(ctx, c, consent.RejectionInfo{
		RejectedBy: api.ConsentRejectedByUSER,
		Reason:     reason,
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
	meta := api.NewRequestMeta(ctx)
	consent, err := s.service.Get(ctx, meta, request.ConsentId)
	if err != nil {
		return nil, err
	}

	resp := newResponse(consent, s.baseURL)
	return api.ConsentV2200JSONResponse(resp), nil
}

func newConsent(
	meta api.RequestMeta,
	req api.CreateConsentRequestV2,
) consent.Consent {
	now := time.Now().UTC()
	c := consent.Consent{
		ID:          consent.ID(),
		Status:      api.ConsentStatusAWAITINGAUTHORISATION,
		UserCPF:     req.Data.LoggedUser.Document.Identification,
		ClientId:    meta.ClientID,
		Permissions: req.Data.Permissions,
		CreatedAt:   now,
		UpdatedAt:   now,
		ExpiresAt:   req.Data.ExpirationDateTime.Time,
	}

	if req.Data.BusinessEntity != nil {
		c.BusinessCNPJ = req.Data.BusinessEntity.Document.Identification
	}

	if req.Data.EndorsementInformation != nil {
		c.EndorsementInfo = &consent.EndorsementInfo{
			PolicyNumber: req.Data.EndorsementInformation.PolicyNumber,
			Type:         req.Data.EndorsementInformation.EndorsementType,
			Description:  req.Data.EndorsementInformation.RequestDescription,
		}
	}

	return c
}

func newResponse(consent consent.Consent, baseURL string) api.ConsentResponseV2 {

	data := api.ConsentDataV2{
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

	if consent.EndorsementInfo != nil {
		data.EndorsementInformation = &api.EndorsementInfo{
			PolicyNumber:       consent.EndorsementInfo.PolicyNumber,
			EndorsementType:    consent.EndorsementInfo.Type,
			RequestDescription: consent.EndorsementInfo.Description,
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
