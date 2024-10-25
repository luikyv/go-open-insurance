package v2

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/resource"
)

type Server struct {
	baseURL string
	service resource.Service
}

func NewServer(baseURL string, service resource.Service) Server {
	return Server{
		baseURL: baseURL,
		service: service,
	}
}

func (s Server) ResourcesV2(
	ctx context.Context,
	request api.ResourcesV2RequestObject,
) (
	api.ResourcesV2ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	pagination := api.NewPagination(request.Params.Page, request.Params.PageSize)
	resources, err := s.service.Resources(ctx, meta, pagination)
	if err != nil {
		return nil, err
	}
	resp := newResourcesResponse(s.baseURL+meta.RequestURI, resources)
	return api.ResourcesV2200JSONResponse(resp), nil
}

func newResourcesResponse(
	requestedURL string,
	page api.Page[api.ResourceData],
) api.ResourcesResponseV2 {
	resp := api.ResourcesResponseV2{
		Data:  page.Records,
		Links: api.PaginatedLinks(requestedURL, page),
		Meta: api.Meta{
			TotalPages:   int32(page.TotalPages),
			TotalRecords: int32(page.TotalRecords),
		},
	}

	return resp
}
