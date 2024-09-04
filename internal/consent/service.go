package consent

import (
	"context"
	"net/http"

	"github.com/luikyv/go-opf/internal/opinerr"
	"github.com/luikyv/go-opf/internal/sec"
	"github.com/luikyv/go-opf/internal/user"
)

type Service struct {
	userService user.Service
	storage     Storage
}

func NewService(userService user.Service, storage Storage) Service {
	return Service{
		userService: userService,
		storage:     storage,
	}
}

func (s Service) Authorize(
	ctx context.Context,
	id string,
	permissions ...Permission,
) error {

	consent, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	if consent.Status != StatusAwaitingAuthorisation {
		return opinerr.New("INVALID_STATUS", http.StatusBadRequest,
			"invalid consent status")
	}

	consent.Status = StatusAuthorised
	consent.Permissions = permissions
	return s.save(ctx, consent)
}

func (s Service) Create(
	ctx context.Context,
	consent Consent,
	meta sec.Meta,
) error {
	if err := s.validate(consent); err != nil {
		return err
	}
	return s.save(ctx, consent)
}

func (s Service) Get(
	ctx context.Context,
	id string,
	meta sec.Meta,
) (
	Consent,
	error,
) {
	consent, err := s.get(ctx, id)
	if err != nil {
		return Consent{}, err
	}

	if consent.ClientId != meta.ClientID {
		return Consent{}, opinerr.New("UNAUTHORIZED", http.StatusUnauthorized,
			"client not authorized")
	}

	return consent, nil
}

func (s Service) Reject(
	ctx context.Context,
	id string,
	meta sec.Meta,
) error {
	consent, err := s.Get(ctx, id, meta)
	if err != nil {
		return err
	}

	if consent.Status == StatusRejected {
		return opinerr.New("INVALID_OPERATION", http.StatusBadRequest,
			"the consent is already rejected")
	}

	consent.Status = StatusRejected
	return s.save(ctx, consent)
}

func (s Service) save(
	ctx context.Context,
	consent Consent,
) error {
	if err := s.storage.Save(ctx, consent); err != nil {
		return opinerr.ErrorInternal
	}
	return nil
}

func (s Service) get(ctx context.Context, id string) (Consent, error) {
	consent, err := s.storage.Get(ctx, id)
	if err != nil {
		return Consent{}, opinerr.New("NOT_FOUND", http.StatusNotFound,
			"could not find the consent")
	}

	if consent.Status != StatusRejected && consent.IsExpired() {
		consent.Status = StatusRejected
		if err := s.save(ctx, consent); err != nil {
			return Consent{}, opinerr.ErrorInternal
		}
	}

	return consent, nil
}

func (s Service) validate(consent Consent) error {
	user, err := s.userService.UserByCPF(consent.UserCPF)
	if err != nil {
		return err
	}

	if consent.BusinessCNPJ != "" && !s.userService.UserBelongsToCompany(user, consent.BusinessCNPJ) {
		return opinerr.New("INVALID_REQUEST", http.StatusBadRequest,
			"the user does not have access to the business entity informed")
	}

	return nil
}
