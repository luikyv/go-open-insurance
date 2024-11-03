package quoteauto

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type ServerV1 struct {
	service Service
}

func NewServerV1(service Service) ServerV1 {
	return ServerV1{
		service: service,
	}
}

func (s ServerV1) CreateQuoteAutoLeadV1(
	ctx context.Context,
	request api.CreateQuoteAutoLeadV1RequestObject,
) (
	api.CreateQuoteAutoLeadV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp, err := s.service.createLead(ctx, meta, *request.Body)
	if err != nil {
		return nil, err
	}

	return api.CreateQuoteAutoLeadV1201JSONResponse(resp), nil
}

func (s ServerV1) RevokeQuoteAutoLeadV1(
	ctx context.Context,
	request api.RevokeQuoteAutoLeadV1RequestObject,
) (
	api.RevokeQuoteAutoLeadV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp, err := s.service.revokeLeadByConsentID(ctx, meta, request.ConsentId)
	if err != nil {
		return nil, err
	}

	return api.RevokeQuoteAutoLeadV1200JSONResponse(resp), nil
}

func (s ServerV1) CreateQuoteAutoV1(
	ctx context.Context,
	request api.CreateQuoteAutoV1RequestObject,
) (
	api.CreateQuoteAutoV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp, err := s.service.createQuote(ctx, meta, *request.Body)
	if err != nil {
		return nil, err
	}

	return api.CreateQuoteAutoV1201JSONResponse(resp), nil
}

func (s ServerV1) QuoteAutoStatusV1(
	ctx context.Context,
	request api.QuoteAutoStatusV1RequestObject,
) (
	api.QuoteAutoStatusV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp, err := s.service.quoteStatus(ctx, meta, request.ConsentId)
	if err != nil {
		return nil, err
	}

	return api.QuoteAutoStatusV1200JSONResponse(resp), nil
}

func (s ServerV1) PatchQuoteAutoV1(
	ctx context.Context,
	request api.PatchQuoteAutoV1RequestObject,
) (
	api.PatchQuoteAutoV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp, err := s.service.patchQuote(ctx, meta, request.ConsentId, *request.Body)
	if err != nil {
		return nil, err
	}

	return api.PatchQuoteAutoV1200JSONResponse(resp), nil
}
