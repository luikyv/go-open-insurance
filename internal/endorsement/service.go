package endorsement

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/consent"
	"github.com/luikyv/go-open-insurance/internal/resource"
)

type Service struct {
	consentService  consent.Service
	resourceService resource.Service
}

func NewService(
	consentService consent.Service,
	resourceService resource.Service,
) Service {
	return Service{
		consentService:  consentService,
		resourceService: resourceService,
	}
}

func (s Service) create(
	ctx context.Context,
	meta api.RequestMeta,
	consentID string,
	req api.CreateEndorsementRequest,
) (
	api.CreateEndorsementResponse,
	error,
) {
	consent, err := s.consentService.FetchAndConsume(ctx, meta, meta.ConsentID)
	if err != nil {
		return api.CreateEndorsementResponse{}, err
	}

	endorsement := newEndorsement(req, consentID)
	if err := s.validate(ctx, meta, endorsement, consent); err != nil {
		return api.CreateEndorsementResponse{}, err
	}

	return newCreateResponse(endorsement), nil
}

func (s Service) validate(
	ctx context.Context,
	meta api.RequestMeta,
	endorsement Endorsement,
	consent consent.Consent,
) error {

	if endorsement.ConsentID != meta.ConsentID {
		return api.NewError("NAO_INFORMADO", http.StatusBadRequest,
			"invalid consent id")
	}

	if meta.Error != nil {
		return api.NewError("NAO_INFORMADO", http.StatusUnprocessableEntity,
			meta.Error.Error())
	}

	info := *consent.Data.EndorsementInformation
	if endorsement.PolicyNumber != info.PolicyNumber {
		return api.NewError("NAO_INFORMADO", http.StatusUnprocessableEntity,
			"policy number not consented")
	}
	if endorsement.Type != info.EndorsementType {
		return api.NewError("NAO_INFORMADO", http.StatusUnprocessableEntity,
			"endorsement type not consented")
	}
	if _, err := s.resourceService.Resource(ctx, meta, endorsement.PolicyNumber); err != nil {
		return api.NewError("NAO_INFORMADO", http.StatusUnprocessableEntity,
			"policy number not found")
	}

	return nil
}

func ID() string {
	return uuid.NewString()
}
