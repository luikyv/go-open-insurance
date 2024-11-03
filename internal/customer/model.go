package customer

import "github.com/luikyv/go-open-insurance/internal/api"

func newPersonalIdentificationsResponse(
	meta api.RequestMeta,
	identifications []api.PersonalIdentificationData,
) api.GetPersonalIdentificationResponse {
	totalPages := 1
	if len(identifications) == 0 {
		totalPages = 0
	}
	resp := api.GetPersonalIdentificationResponse{
		Data: identifications,
		Links: api.Links{
			Self: meta.RequestURL(),
		},
		Meta: api.Meta{
			TotalPages:   int32(totalPages),
			TotalRecords: int32(len(identifications)),
		},
	}

	return resp
}

func newPersonalQualificationsResponse(
	meta api.RequestMeta,
	qualifications []api.PersonalQualificationData,
) api.GetPersonalQualificationResponse {
	totalPages := 1
	if len(qualifications) == 0 {
		totalPages = 0
	}
	resp := api.GetPersonalQualificationResponse{
		Data: qualifications,
		Links: api.Links{
			Self: meta.RequestURL(),
		},
		Meta: api.Meta{
			TotalPages:   int32(totalPages),
			TotalRecords: int32(len(qualifications)),
		},
	}

	return resp
}

func newPersonalComplimentaryInfoResponse(
	meta api.RequestMeta,
	infos []api.PersonalComplimentaryInfoData,
) api.GetPersonalComplimentaryInfoResponse {
	totalPages := 1
	if len(infos) == 0 {
		totalPages = 0
	}
	resp := api.GetPersonalComplimentaryInfoResponse{
		Data: infos,
		Links: api.Links{
			Self: meta.RequestURL(),
		},
		Meta: api.Meta{
			TotalPages:   int32(totalPages),
			TotalRecords: int32(len(infos)),
		},
	}

	return resp
}
