package consent

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type ServerV2 struct {
	service Service
}

func NewServerV2(
	service Service,
) ServerV2 {
	return ServerV2{
		service: service,
	}
}

func (s ServerV2) CreateConsentV2(
	ctx context.Context,
	request api.CreateConsentV2RequestObject,
) (
	api.CreateConsentV2ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp, err := s.service.create(ctx, meta, *request.Body)
	if err != nil {
		return nil, err
	}

	return api.CreateConsentV2201JSONResponse(resp), nil
}

func (s ServerV2) DeleteConsentV2(
	ctx context.Context,
	request api.DeleteConsentV2RequestObject,
) (
	api.DeleteConsentV2ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	if err := s.service.delete(ctx, meta, request.ConsentId); err != nil {
		return nil, err
	}

	return api.DeleteConsentV2204Response{}, nil
}

func (s ServerV2) ConsentV2(
	ctx context.Context,
	request api.ConsentV2RequestObject,
) (
	api.ConsentV2ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp, err := s.service.fetch(ctx, meta, request.ConsentId)
	if err != nil {
		return nil, err
	}

	return api.ConsentV2200JSONResponse(resp), nil
}
