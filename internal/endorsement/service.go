package endorsement

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/luikyv/go-open-insurance/internal/consent"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
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

func (s Service) Create(
	ctx context.Context,
	sub string,
	endorsement Endorsement,
) error {
	consent, err := s.consentService.GetAndConsume(ctx, endorsement.ConsentID)
	if err != nil {
		return err
	}

	if err := s.validate(ctx, sub, endorsement, consent); err != nil {
		return err
	}

	return nil
}

func (s Service) validate(
	ctx context.Context,
	sub string,
	endorsement Endorsement,
	consent consent.Consent,
) error {
	info := *consent.EndorsementInfo

	if endorsement.PolicyNumber != info.PolicyNumber {
		return opinerr.New("NAO_INFORMADO", http.StatusUnprocessableEntity,
			"policy number not consented")
	}

	if endorsement.Type != info.Type {
		return opinerr.New("NAO_INFORMADO", http.StatusUnprocessableEntity,
			"endorsement type not consented")
	}

	if _, err := s.resourceService.Get(ctx, sub, endorsement.PolicyNumber); err != nil {
		return opinerr.New("NAO_INFORMADO", http.StatusUnprocessableEntity,
			"policy number not found")
	}

	return nil
}

func ID() string {
	return uuid.NewString()
}
