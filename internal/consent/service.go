package consent

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/user"
)

type Service struct {
	storage     Storage
	userService user.Service
}

func NewService(storage Storage, userService user.Service) Service {
	return Service{
		storage:     storage,
		userService: userService,
	}
}

func (s Service) Authorize(
	ctx context.Context,
	id string,
	permissions ...api.ConsentPermission,
) error {

	api.Logger(ctx).Debug("trying to authorize consent",
		slog.String("consent_id", id))
	consent, err := s.fetchAndModify(ctx, id)
	if err != nil {
		return err
	}

	if !consent.IsAwaitingAuthorization() {
		api.Logger(ctx).Debug("cannot authorize a consent that is not awaiting authorization",
			slog.String("consent_id", id), slog.Any("status", consent.Status))
		return api.NewError("INVALID_STATUS", http.StatusBadRequest,
			"invalid consent status")
	}

	api.Logger(ctx).Info("authorizing consent",
		slog.String("consent_id", id))
	consent.Status = api.ConsentStatusAUTHORISED
	consent.Permissions = permissions
	return s.save(ctx, consent)
}

func (s Service) Fetch(
	ctx context.Context,
	meta api.RequestMeta,
	id string,
) (
	Consent,
	error,
) {
	consent, err := s.fetchAndModify(ctx, id)
	if err != nil {
		return Consent{}, err
	}

	if meta.ClientID != consent.ClientId {
		api.Logger(ctx).Debug("client not allowed to fetch the consent")
		return Consent{}, api.NewError("UNAUTHORIZED", http.StatusForbidden,
			"client not authorized to perform this operation")
	}

	return consent, nil
}

func (s Service) FetchAndConsume(
	ctx context.Context,
	meta api.RequestMeta,
	id string,
) (
	Consent,
	error,
) {
	consent, err := s.Fetch(ctx, meta, id)
	if err != nil {
		return Consent{}, err
	}
	if err := s.consume(ctx, consent); err != nil {
		return Consent{}, err
	}

	return consent, nil
}

func (s Service) Reject(
	ctx context.Context,
	meta api.RequestMeta,
	id string,
	info RejectionInfo,
) error {
	consent, err := s.Fetch(ctx, meta, id)
	if err != nil {
		return err
	}

	return s.reject(ctx, consent, info)
}

// Verify checks if the consent with the given ID is authorized
// and has the required permissions.
func (s Service) Verify(
	ctx context.Context,
	meta api.RequestMeta,
	id string,
	permissions ...api.ConsentPermission,
) error {
	consent, err := s.Fetch(ctx, meta, id)
	if err != nil {
		return err
	}

	if !consent.IsAuthorized() {
		return api.NewError("INVALID_STATUS", http.StatusBadRequest,
			"consent is not authorized")
	}

	if !consent.HasPermissions(permissions) {
		return api.NewError("INVALID_PERMISSIONS", http.StatusBadRequest,
			"consent missing permissions")
	}

	return nil
}

func (s Service) create(
	ctx context.Context,
	meta api.RequestMeta,
	req api.CreateConsentRequest,
) (
	api.ConsentResponse,
	error,
) {
	consent := newConsent(meta, req)
	if err := s.validate(ctx, meta, consent); err != nil {
		return api.ConsentResponse{}, err
	}

	api.Logger(ctx).Info("creating consent", slog.String("consent_id", consent.ID))
	if err := s.save(ctx, consent); err != nil {
		return api.ConsentResponse{}, err
	}

	return newResponse(meta, consent), nil
}

func (s Service) consume(
	ctx context.Context,
	consent Consent,
) error {
	if !consent.IsAuthorized() {
		return api.NewError("INVALID_OPERATION", http.StatusBadRequest,
			"cannot consume a consent that is not authorized")
	}

	consent.Status = api.ConsentStatusCONSUMED
	return s.save(ctx, consent)
}

func (s Service) delete(
	ctx context.Context,
	meta api.RequestMeta,
	id string,
) error {
	c, err := s.Fetch(ctx, meta, id)
	if err != nil {
		return err
	}

	reason := api.ConsentRejectedReasonCodeCUSTOMERMANUALLYREJECTED
	if c.IsAuthorized() {
		reason = api.ConsentRejectedReasonCodeCUSTOMERMANUALLYREVOKED
	}

	if err := s.reject(ctx, c, RejectionInfo{
		RejectedBy: api.ConsentRejectedByUSER,
		Reason:     reason,
	}); err != nil {
		return err
	}

	return s.save(ctx, c)
}

func (s Service) reject(
	ctx context.Context,
	consent Consent,
	info RejectionInfo,
) error {
	if consent.Status == api.ConsentStatusREJECTED {
		return api.NewError("INVALID_OPERATION", http.StatusBadRequest,
			"the consent is already rejected")
	}

	consent.Status = api.ConsentStatusREJECTED
	consent.RejectionInfo = &info
	return s.save(ctx, consent)
}

func (s Service) save(
	ctx context.Context,
	consent Consent,
) error {
	if err := s.storage.save(ctx, consent); err != nil {
		api.Logger(ctx).Error("could not save the consent", slog.Any("error", err))
		return api.ErrInternal
	}
	return nil
}

func (s Service) fetch(
	ctx context.Context,
	meta api.RequestMeta,
	id string,
) (
	api.ConsentResponse,
	error,
) {
	consent, err := s.Fetch(ctx, meta, id)
	if err != nil {
		return api.ConsentResponse{}, err
	}

	return newResponse(meta, consent), nil
}

func (s Service) fetchAndModify(ctx context.Context, id string) (Consent, error) {
	consent, err := s.storage.fetch(ctx, id)
	if err != nil {
		api.Logger(ctx).Debug("could not find the consent", slog.Any("error", err))
		return Consent{}, api.NewError("NOT_FOUND", http.StatusNotFound,
			"could not find the consent")
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
			RejectedBy: api.ConsentRejectedByASPSP,
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
func (s Service) validate(ctx context.Context, meta api.RequestMeta, consent Consent) error {
	if meta.Error != nil {
		return api.NewError("NAO_INFORMADO", http.StatusBadRequest,
			meta.Error.Error())
	}

	if err := validate(ctx, consent); err != nil {
		api.Logger(ctx).Debug("the consent is not valid", slog.Any("error", err))
		return err
	}

	return nil
}
