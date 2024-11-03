package resource

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type ServerV2 struct {
	service Service
}

func NewServerV2(service Service) ServerV2 {
	return ServerV2{
		service: service,
	}
}

func (s ServerV2) ResourcesV2(
	ctx context.Context,
	request api.ResourcesV2RequestObject,
) (
	api.ResourcesV2ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	pagination := api.NewPagination(request.Params.Page, request.Params.PageSize)
	resp, err := s.service.resources(ctx, meta, pagination)
	if err != nil {
		return nil, err
	}

	return api.ResourcesV2200JSONResponse(resp), nil
}
