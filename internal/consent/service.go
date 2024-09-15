package consent

import (
	"context"
	"log/slog"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
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
	permissions ...api.ConsentPermission,
) error {

	api.Logger(ctx).Debug("trying to authorize consent",
		slog.String("consent_id", id))
	consent, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	if consent.Status != api.ConsentStatusAWAITINGAUTHORISATION {
		api.Logger(ctx).Debug("cannot authorize a consent that is not awaiting authorization",
			slog.String("consent_id", id), slog.Any("status", consent.Status))
		return errInvalidStatus
	}

	api.Logger(ctx).Info("authorizing consent",
		slog.String("consent_id", id))
	consent.Status = api.ConsentStatusAUTHORISED
	consent.Permissions = permissions
	return s.save(ctx, consent)
}

func (s Service) Create(
	ctx context.Context,
	consent Consent,
) error {
	if err := s.validate(ctx, consent); err != nil {
		return err
	}

	api.Logger(ctx).Info("creating consent", slog.String("consent_id", consent.ID))
	return s.save(ctx, consent)
}

func (s Service) Get(
	ctx context.Context,
	id string,
) (
	Consent,
	error,
) {
	consent, err := s.get(ctx, id)
	if err != nil {
		return Consent{}, err
	}

	clientID := ctx.Value(api.CtxKeyClientID)
	if clientID != consent.ClientId {
		api.Logger(ctx).Debug("client not allowed to fetch the consent")
		return Consent{}, errClientNotAuthorized
	}

	return consent, nil
}

func (s Service) Reject(
	ctx context.Context,
	id string,
	info RejectionInfo,
) error {
	consent, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if consent.Status == api.ConsentStatusREJECTED {
		return errAlreadyRejected
	}

	consent.Status = api.ConsentStatusREJECTED
	consent.RejectionInfo = &info
	return s.save(ctx, consent)
}

func (s Service) VerifyPermissions(
	ctx context.Context,
	id string,
	permissions ...api.ConsentPermission,
) error {
	consent, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if !consent.IsAuthorized() {
		return errInvalidStatus
	}

	if !consent.HasPermissions(permissions) {
		return errInvalidPermissions
	}
	return nil
}

func (s Service) save(
	ctx context.Context,
	consent Consent,
) error {
	if err := s.storage.Save(ctx, consent); err != nil {
		api.Logger(ctx).Error("could not save the consent", slog.Any("error", err))
		return opinerr.ErrInternal
	}
	return nil
}

func (s Service) get(ctx context.Context, id string) (Consent, error) {
	consent, err := s.storage.Get(ctx, id)
	if err != nil {
		api.Logger(ctx).Debug("could not find the consent", slog.Any("error", err))
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
		api.Logger(ctx).Debug("consent awaiting authorization for too long, moving to rejected")
		consent.Status = api.ConsentStatusREJECTED
		consent.RejectionInfo = &RejectionInfo{
			RejectedBy: api.ConsentRejectedByUSER,
			Reason:     api.ConsentRejectedReasonCodeCONSENTEXPIRED,
		}
		consentWasModified = true
	}

	// Reject the consent if it reached the expiration.
	if consent.IsExpired() {
		api.Logger(ctx).Debug("consent reached expiration, moving to rejected")
		consent.Status = api.ConsentStatusREJECTED
		consent.RejectionInfo = &RejectionInfo{
			RejectedBy: api.ConsentRejectedByUSER,
			Reason:     api.ConsentRejectedReasonCodeCONSENTMAXDATEREACHED,
		}
		consentWasModified = true
	}

	if consentWasModified {
		api.Logger(ctx).Debug("the consent was modified")
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
		api.Logger(ctx).Debug("the consent is not valid", slog.Any("error", err))
		return err
	}

	return nil
}
