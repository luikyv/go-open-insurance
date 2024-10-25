package v1

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/capitalizationtitle"
)

type Server struct {
	baseURL string
	service capitalizationtitle.Service
}

func NewServer(
	baseURL string,
	service capitalizationtitle.Service,
) Server {
	return Server{
		baseURL: baseURL,
		service: service,
	}
}

func (s Server) CapitalizationTitlePlans(
	ctx context.Context,
	request api.CapitalizationTitlePlansV1RequestObject,
) (
	api.CapitalizationTitlePlansV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	pagination := api.NewPagination(request.Params.Page, request.Params.PageSize)

	plans := s.service.Plans(meta, pagination)
	resp := newPlansResponse(s.baseURL+meta.RequestURI, plans)
	return api.CapitalizationTitlePlansV1200JSONResponse(resp), nil
}

func (s Server) CapitalizationTitleEvents(
	ctx context.Context,
	request api.CapitalizationTitleEventsV1RequestObject,
) (
	api.CapitalizationTitleEventsV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	pagination := api.NewPagination(request.Params.Page, request.Params.PageSize)

	events, err := s.service.PlanEvents(meta, request.PlanId, pagination)
	if err != nil {
		return nil, err
	}

	resp := newPlanEventsResponse(s.baseURL+meta.RequestURI, events)
	return api.CapitalizationTitleEventsV1200JSONResponse(resp), nil
}

func (s Server) CapitalizationTitlePlanInfo(
	ctx context.Context,
	request api.CapitalizationTitlePlanInfoV1RequestObject,
) (
	api.CapitalizationTitlePlanInfoV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	info, err := s.service.PlanInfo(meta, request.PlanId)
	if err != nil {
		return nil, err
	}

	resp := newPlanIfnoResponse(s.baseURL+meta.RequestURI, info)
	return api.CapitalizationTitlePlanInfoV1200JSONResponse(resp), nil
}

func (s Server) CapitalizationTitleSettlements(
	ctx context.Context,
	request api.CapitalizationTitleSettlementsV1RequestObject,
) (
	api.CapitalizationTitleSettlementsV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	pagination := api.NewPagination(request.Params.Page, request.Params.PageSize)

	settlements, err := s.service.PlanSettlements(meta, request.PlanId, pagination)
	if err != nil {
		return nil, err
	}

	resp := newPlanSettlementsResponse(s.baseURL+meta.RequestURI, settlements)
	return api.CapitalizationTitleSettlementsV1200JSONResponse(resp), nil
}

func newPlansResponse(
	requestedURL string,
	page api.Page[api.CapitalizationTitlePlanData],
) api.CapitalizationTitlePlansResponseV1 {
	return api.CapitalizationTitlePlansResponseV1{
		Data:  page.Records,
		Links: api.PaginatedLinks(requestedURL, page),
		Meta: api.Meta{
			TotalPages:   int32(page.TotalPages),
			TotalRecords: int32(page.TotalRecords),
		},
	}
}

func newPlanEventsResponse(
	requestedURL string,
	page api.Page[api.CapitalizationTitleEvent],
) api.CapitalizationTitleEventsResponseV1 {
	return api.CapitalizationTitleEventsResponseV1{
		Data:  page.Records,
		Links: api.PaginatedLinks(requestedURL, page),
		Meta: api.Meta{
			TotalPages:   int32(page.TotalPages),
			TotalRecords: int32(page.TotalRecords),
		},
	}
}

func newPlanIfnoResponse(
	requestedURL string,
	info api.CapitalizationTitlePlanInfo,
) api.CapitalizationTitlePlanInfoResponseV1 {
	return api.CapitalizationTitlePlanInfoResponseV1{
		Data: info,
		Links: api.Links{
			Self: requestedURL,
		},
		Meta: api.Meta{
			TotalPages:   1,
			TotalRecords: 1,
		},
	}
}

func newPlanSettlementsResponse(
	requestedURL string,
	page api.Page[api.CapitalizationTitleSettlement],
) api.CapitalizationTitleSettlementsResponseV1 {
	return api.CapitalizationTitleSettlementsResponseV1{
		Data:  page.Records,
		Links: api.PaginatedLinks(requestedURL, page),
		Meta: api.Meta{
			TotalPages:   int32(page.TotalPages),
			TotalRecords: int32(page.TotalRecords),
		},
	}
}
