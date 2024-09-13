package consent

import (
	"context"
	"log/slog"

	"github.com/luikyv/go-open-insurance/internal/log"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
	"github.com/luikyv/go-open-insurance/internal/sec"
	"github.com/luikyv/go-open-insurance/internal/user"
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

	log.FromCtx(ctx).Debug("trying to authorize consent",
		slog.String("consent_id", id))
	consent, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	if consent.Status != StatusAwaitingAuthorisation {
		log.FromCtx(ctx).Debug("cannot authorize a consent that is not awaiting authorization",
			slog.String("consent_id", id), slog.Any("status", consent.Status))
		return errInvalidStatus
	}

	log.FromCtx(ctx).Info("authorizing consent",
		slog.String("consent_id", id))
	consent.Status = StatusAuthorised
	consent.Permissions = permissions
	return s.save(ctx, consent)
}

func (s Service) Create(
	ctx context.Context,
	meta sec.Meta,
	consent Consent,
) error {
	if err := s.validate(ctx, consent); err != nil {
		return err
	}

	log.FromCtx(ctx).Info("creating consent", slog.String("consent_id", consent.ID))
	return s.save(ctx, consent)
}

func (s Service) Get(
	ctx context.Context,
	meta sec.Meta,
	id string,
) (
	Consent,
	error,
) {
	consent, err := s.get(ctx, id)
	if err != nil {
		return Consent{}, err
	}

	if consent.ClientId != meta.ClientID {
		log.FromCtx(ctx).Debug("client not allowed to fetch the consent",
			slog.String("client_id", meta.ClientID))
		return Consent{}, errClientNotAuthorized
	}

	return consent, nil
}

func (s Service) Reject(
	ctx context.Context,
	meta sec.Meta,
	id string,
	info RejectionInfo,
) error {
	consent, err := s.Get(ctx, meta, id)
	if err != nil {
		return err
	}

	if consent.Status == StatusRejected {
		return errAlreadyRejected
	}

	consent.Status = StatusRejected
	consent.RejectionInfo = &info
	return s.save(ctx, consent)
}

func (s Service) save(
	ctx context.Context,
	consent Consent,
) error {
	if err := s.storage.Save(ctx, consent); err != nil {
		log.FromCtx(ctx).Error("could not save the consent", slog.Any("error", err))
		return opinerr.ErrInternal
	}
	return nil
}

func (s Service) get(ctx context.Context, id string) (Consent, error) {
	consent, err := s.storage.Get(ctx, id)
	if err != nil {
		log.FromCtx(ctx).Debug("could not find the consent", slog.Any("error", err))
		return Consent{}, errNotFound
	}

	if err := s.modify(ctx, &consent); err != nil {
		return Consent{}, err
	}

	return consent, nil
}

// modify will evaluated the consent information and modify it to be compliant.
func (s Service) modify(ctx context.Context, consent *Consent) error {
	consentWasModified := false

	// Reject the consent if the time awaiting the user authorization has elapsed.
	if consent.HasAuthExpired() {
		log.FromCtx(ctx).Debug("consent awaiting authorization for too long, moving to rejected")
		consent.Status = StatusRejected
		consent.RejectionInfo = &RejectionInfo{
			RejectedBy: RejectedByUser,
			Reason:     RejectionReasonConsentExpired,
		}
		consentWasModified = true
	}

	// Reject the consent if it reached the expiration.
	if consent.IsExpired() {
		log.FromCtx(ctx).Debug("consent reached expiration, moving to rejected")
		consent.Status = StatusRejected
		consent.RejectionInfo = &RejectionInfo{
			RejectedBy: RejectedByUser,
			Reason:     RejectionReasonConsentMaxDateReached,
		}
		consentWasModified = true
	}

	if consentWasModified {
		log.FromCtx(ctx).Debug("the consent was modified")
		if err := s.save(ctx, *consent); err != nil {
			return err
		}
	}

	return nil
}

// validate validates the consent information.
// This is intended to be used before the consent is created to make sure the
// information is compliant.
func (s Service) validate(ctx context.Context, consent Consent) error {
	if err := validate(ctx, consent); err != nil {
		log.FromCtx(ctx).Debug("the consent is not valid", slog.Any("error", err))
	}

	return nil
}
