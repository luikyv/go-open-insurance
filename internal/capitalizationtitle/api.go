package capitalizationtitle

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type ServerV1 struct {
	service Service
}

func NewServerV1(
	service Service,
) ServerV1 {
	return ServerV1{
		service: service,
	}
}

func (s ServerV1) CapitalizationTitlePlansV1(
	ctx context.Context,
	request api.CapitalizationTitlePlansV1RequestObject,
) (
	api.CapitalizationTitlePlansV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	pagination := api.NewPagination(request.Params.Page, request.Params.PageSize)

	resp := s.service.plans(meta, pagination)
	return api.CapitalizationTitlePlansV1200JSONResponse(resp), nil
}

func (s ServerV1) CapitalizationTitleEventsV1(
	ctx context.Context,
	request api.CapitalizationTitleEventsV1RequestObject,
) (
	api.CapitalizationTitleEventsV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	pagination := api.NewPagination(request.Params.Page, request.Params.PageSize)

	resp, err := s.service.planEvents(meta, request.PlanId, pagination)
	if err != nil {
		return nil, err
	}

	return api.CapitalizationTitleEventsV1200JSONResponse(resp), nil
}

func (s ServerV1) CapitalizationTitlePlanInfoV1(
	ctx context.Context,
	request api.CapitalizationTitlePlanInfoV1RequestObject,
) (
	api.CapitalizationTitlePlanInfoV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp, err := s.service.planInfo(meta, request.PlanId)
	if err != nil {
		return nil, err
	}

	return api.CapitalizationTitlePlanInfoV1200JSONResponse(resp), nil
}

func (s ServerV1) CapitalizationTitleSettlementsV1(
	ctx context.Context,
	request api.CapitalizationTitleSettlementsV1RequestObject,
) (
	api.CapitalizationTitleSettlementsV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	pagination := api.NewPagination(request.Params.Page, request.Params.PageSize)

	resp, err := s.service.planSettlements(meta, request.PlanId, pagination)
	if err != nil {
		return nil, err
	}

	return api.CapitalizationTitleSettlementsV1200JSONResponse(resp), nil
}
