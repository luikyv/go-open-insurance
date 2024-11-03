package quoteauto

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/webhook"
)

type Service struct {
	storage        Storage
	webhookService webhook.Service
}

func NewService(
	storage Storage,
	webhookService webhook.Service,
) Service {
	return Service{
		storage:        storage,
		webhookService: webhookService,
	}
}

func (s Service) createLead(
	ctx context.Context,
	meta api.RequestMeta,
	req api.CreateQuoteAutoLeadRequest,
) (
	api.CreateQuoteLeadResponse,
	error,
) {
	lead := newLead(req)
	if err := s.validateLead(ctx, meta, lead); err != nil {
		return api.CreateQuoteLeadResponse{}, err
	}

	if err := s.saveLead(ctx, lead); err != nil {
		return api.CreateQuoteLeadResponse{}, err
	}

	return newLeadCreateResponse(meta, lead), nil
}

func (s Service) revokeLeadByConsentID(
	ctx context.Context,
	meta api.RequestMeta,
	consentID string,
) (
	api.RevokeQuoteLeadResponse,
	error,
) {
	lead, err := s.storage.fetchLeadByConsentID(ctx, consentID)
	if err != nil {
		return api.RevokeQuoteLeadResponse{}, api.NewError("NOT_FOUND", http.StatusNotFound,
			"the auto quote lead was not found for the informed consent id")
	}

	if err := s.validateLeadRevocation(ctx, meta); err != nil {
		return api.RevokeQuoteLeadResponse{}, err
	}

	if err := s.revokeLead(ctx, &lead); err != nil {
		return api.RevokeQuoteLeadResponse{}, err
	}

	return newRevokeLeadResponse(), nil
}

func (s Service) revokeLead(ctx context.Context, lead *Lead) error {
	lead.Status = api.QuoteStatusCANC
	lead.StatusUpdateDateTime = time.Now().UTC()
	return s.saveLead(ctx, *lead)
}

func (s Service) saveLead(
	ctx context.Context,
	lead Lead,
) error {
	if err := s.storage.saveLead(ctx, lead); err != nil {
		api.Logger(ctx).Error("could not save auto quote lead",
			slog.String("error", err.Error()))
		return api.ErrInternal
	}

	return nil
}

func (s Service) createQuote(
	ctx context.Context,
	meta api.RequestMeta,
	req api.CreateQuoteAutoRequest,
) (
	api.CreateQuoteResponse,
	error,
) {
	if err := s.validateCreateQuoteRequest(ctx, meta, req); err != nil {
		return api.CreateQuoteResponse{}, err
	}

	quote := newQuote(req)
	if err := s.saveQuote(ctx, &quote); err != nil {
		return api.CreateQuoteResponse{}, err
	}

	return newCreateQuoteResponse(meta, quote), nil
}

func (s Service) quoteStatus(
	ctx context.Context,
	meta api.RequestMeta,
	consentID string,
) (
	api.GetQuoteAutoStatusResponse,
	error,
) {
	quote, err := s.quoteByConsentID(ctx, consentID)
	if err != nil {
		return api.GetQuoteAutoStatusResponse{}, err
	}

	if err := s.modifyQuote(ctx, meta, &quote); err != nil {
		return api.GetQuoteAutoStatusResponse{}, err
	}

	return newGetQuoteAutoStatusResponse(meta, quote), nil
}

func (s Service) modifyQuote(
	ctx context.Context,
	meta api.RequestMeta,
	quote *Quote,
) error {
	quoteWasModified := false

	// The conformance suite rejection test sends the term end date prior to the
	// term start date.
	if quote.Status == api.QuoteStatusEVAL &&
		quote.Data.QuoteData.TermEndDate.Before(quote.Data.QuoteData.TermStartDate.Time) {
		quote.Status = api.QuoteStatusRJCT
		quoteWasModified = true
	}

	if quote.Status == api.QuoteStatusEVAL {
		quote.Status = api.QuoteStatusACPT
		quoteWasModified = true
	}

	if quote.Status == api.QuoteStatusRCVD {
		quote.Status = api.QuoteStatusEVAL
		quoteWasModified = true
	}

	if !quoteWasModified {
		return nil
	}

	s.webhookService.Notify(
		ctx,
		meta.ClientID,
		fmt.Sprintf("/quote/v1/request/%s/quote-status", quote.ConsentID),
	)
	return s.saveQuote(ctx, quote)
}

func (s Service) patchQuote(
	ctx context.Context,
	meta api.RequestMeta,
	consentID string,
	req api.PatchQuoteRequest,
) (
	api.PatchQuoteResponse,
	error,
) {
	if meta.Error != nil {
		return api.PatchQuoteResponse{}, meta.Error
	}

	quote, err := s.quoteByConsentID(ctx, consentID)
	if err != nil {
		return api.PatchQuoteResponse{}, err
	}

	if req.Data.Status == api.PatchQuoteRequestDataStatusACKN {
		err = s.acknowledgeQuote(ctx, &quote)
	} else {
		err = s.cancelQuote(ctx, &quote)
	}
	if err != nil {
		return api.PatchQuoteResponse{}, err
	}

	return newPatchQuoteResponse(meta, quote), nil
}

func (s Service) acknowledgeQuote(
	ctx context.Context,
	quote *Quote,
) error {

	if quote.Status != api.QuoteStatusACPT {
		return api.NewError("NAO_INFORMADO", http.StatusUnprocessableEntity,
			"the quote is not accepted")
	}

	quote.Status = api.QuoteStatusACKN
	return s.saveQuote(ctx, quote)
}

func (s Service) cancelQuote(
	ctx context.Context,
	quote *Quote,
) error {
	quote.Status = api.QuoteStatusCANC
	return s.saveQuote(ctx, quote)
}

func (s Service) quoteByConsentID(
	ctx context.Context,
	id string,
) (
	Quote,
	error,
) {
	quote, err := s.storage.fetchQuoteByConsentID(ctx, id)
	if err != nil {
		return Quote{}, api.NewError("NOT_FOUND", http.StatusNotFound,
			fmt.Sprintf("could not auto quote for consent id %s", id))
	}

	return quote, nil
}

func (s Service) validateCreateQuoteRequest(
	_ context.Context,
	meta api.RequestMeta,
	_ api.CreateQuoteAutoRequest,
) error {
	if meta.Error != nil {
		return api.NewError("NAO_INFORMADO", http.StatusUnprocessableEntity,
			meta.Error.Error())
	}

	return nil
}

func (s Service) validateLead(
	_ context.Context,
	meta api.RequestMeta,
	_ Lead,
) error {
	if meta.Error != nil {
		return api.NewError("NAO_INFORMADO", http.StatusUnprocessableEntity,
			meta.Error.Error())
	}

	return nil
}

func (s Service) validateLeadRevocation(
	_ context.Context,
	meta api.RequestMeta,
) error {
	if meta.Error != nil {
		return api.NewError("NAO_INFORMADO", http.StatusUnprocessableEntity,
			meta.Error.Error())
	}

	return nil
}

func (s Service) saveQuote(
	ctx context.Context,
	quote *Quote,
) error {

	quote.StatusUpdateDateTime = time.Now().UTC()
	if err := s.storage.saveQuote(ctx, *quote); err != nil {
		api.Logger(ctx).Error("could not save auto quote",
			slog.String("error", err.Error()))
		return api.ErrInternal
	}
	return nil
}
