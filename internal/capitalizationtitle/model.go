package capitalizationtitle

import "github.com/luikyv/go-open-insurance/internal/api"

func newPlansResponse(
	meta api.RequestMeta,
	page api.Page[api.CapitalizationTitlePlanData],
) api.GetCapitalizationTitlePlansResponse {
	return api.GetCapitalizationTitlePlansResponse{
		Data:  page.Records,
		Links: api.PaginatedLinks(meta.RequestURL(), page),
		Meta: api.Meta{
			TotalPages:   int32(page.TotalPages),
			TotalRecords: int32(page.TotalRecords),
		},
	}
}

func newPlanEventsResponse(
	meta api.RequestMeta,
	page api.Page[api.CapitalizationTitleEvent],
) api.GetCapitalizationTitleEventsResponse {
	return api.GetCapitalizationTitleEventsResponse{
		Data:  page.Records,
		Links: api.PaginatedLinks(meta.RequestURL(), page),
		Meta: api.Meta{
			TotalPages:   int32(page.TotalPages),
			TotalRecords: int32(page.TotalRecords),
		},
	}
}

func newPlanInfoResponse(
	meta api.RequestMeta,
	info api.CapitalizationTitlePlanInfo,
) api.GetCapitalizationTitlePlanInfoResponse {
	return api.GetCapitalizationTitlePlanInfoResponse{
		Data: info,
		Links: api.Links{
			Self: meta.RequestURL(),
		},
		Meta: api.Meta{
			TotalPages:   1,
			TotalRecords: 1,
		},
	}
}

func newPlanSettlementsResponse(
	meta api.RequestMeta,
	page api.Page[api.CapitalizationTitleSettlement],
) api.GetCapitalizationTitleSettlementsResponse {
	return api.GetCapitalizationTitleSettlementsResponse{
		Data:  page.Records,
		Links: api.PaginatedLinks(meta.RequestURL(), page),
		Meta: api.Meta{
			TotalPages:   int32(page.TotalPages),
			TotalRecords: int32(page.TotalRecords),
		},
	}
}
