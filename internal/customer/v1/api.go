package customersv1

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/customer"
)

type Server struct {
	baseURL string
	service customer.Service
}

func NewServer(baseURL string, service customer.Service) Server {
	return Server{
		baseURL: baseURL,
		service: service,
	}
}

func (s Server) PersonalIdentificationsV1(
	ctx context.Context,
	request api.PersonalIdentificationsV1RequestObject,
) (
	api.PersonalIdentificationsV1ResponseObject,
	error,
) {
	sub := ctx.Value(api.CtxKeySubject).(string)
	identifications := s.service.PersonalIdentifications(sub)
	resp := newPersonalIdentificationsResponse(s.baseURL, identifications)
	return api.PersonalIdentificationsV1200JSONResponse(resp), nil
}

func (s Server) PersonalQualificationsV1(
	ctx context.Context,
	request api.PersonalQualificationsV1RequestObject,
) (
	api.PersonalQualificationsV1ResponseObject,
	error,
) {
	sub := ctx.Value(api.CtxKeySubject).(string)
	qualifications := s.service.PersonalQualifications(sub)
	resp := newPersonalQualificationsResponse(s.baseURL, qualifications)
	return api.PersonalQualificationsV1200JSONResponse(resp), nil
}

func (s Server) PersonalComplimentaryInfoV1(
	ctx context.Context,
	request api.PersonalComplimentaryInfoV1RequestObject,
) (
	api.PersonalComplimentaryInfoV1ResponseObject,
	error,
) {
	sub := ctx.Value(api.CtxKeySubject).(string)
	infos := s.service.PersonalComplimentaryInfos(sub)
	resp := newPersonalComplimentaryInfoResponse(s.baseURL, infos)
	return api.PersonalComplimentaryInfoV1200JSONResponse(resp), nil
}

func newPersonalIdentificationsResponse(
	baseURL string,
	identifications []api.PersonalIdentificationData,
) api.PersonalIdentificationResponseV1 {
	totalPages := 1
	if len(identifications) == 0 {
		totalPages = 0
	}
	resp := api.PersonalIdentificationResponseV1{
		Data: identifications,
		Links: api.Links{
			Self: baseURL + "/customers/v1/personal/identifications",
		},
		Meta: api.Meta{
			TotalPages:   int32(totalPages),
			TotalRecords: int32(len(identifications)),
		},
	}

	return resp
}

func newPersonalQualificationsResponse(
	baseURL string,
	qualifications []api.PersonalQualificationData,
) api.PersonalQualificationResponseV1 {
	totalPages := 1
	if len(qualifications) == 0 {
		totalPages = 0
	}
	resp := api.PersonalQualificationResponseV1{
		Data: qualifications,
		Links: api.Links{
			Self: baseURL + "/customers/v1/personal/qualifications",
		},
		Meta: api.Meta{
			TotalPages:   int32(totalPages),
			TotalRecords: int32(len(qualifications)),
		},
	}

	return resp
}

func newPersonalComplimentaryInfoResponse(
	baseURL string,
	infos []api.PersonalComplimentaryInfoData,
) api.PersonalComplimentaryInfoResponseV1 {
	totalPages := 1
	if len(infos) == 0 {
		totalPages = 0
	}
	resp := api.PersonalComplimentaryInfoResponseV1{
		Data: infos,
		Links: api.Links{
			Self: baseURL + "/customers/v1/personal/complimentary-information",
		},
		Meta: api.Meta{
			TotalPages:   int32(totalPages),
			TotalRecords: int32(len(infos)),
		},
	}

	return resp
}
