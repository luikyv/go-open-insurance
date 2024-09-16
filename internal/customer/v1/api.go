package customersv1

import (
	"context"
	"fmt"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type Server struct {
	baseURL string
	service Service
}

func NewServer(baseURL string, service Service) Server {
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

type Service interface {
	AddPersonalIdentifications(
		sub string,
		identifications []api.PersonalIdentificationDataV1,
	)
	PersonalIdentifications(sub string) []api.PersonalIdentificationDataV1
}

type service struct {
	// personalIdentificationsMap maps users to their identifications.
	personalIdentificationsMap map[string][]api.PersonalIdentificationDataV1
}

func NewService() Service {
	return &service{
		personalIdentificationsMap: map[string][]api.PersonalIdentificationDataV1{},
	}
}

func (s *service) AddPersonalIdentifications(
	sub string,
	identifications []api.PersonalIdentificationDataV1,
) {
	s.personalIdentificationsMap[sub] = identifications
}

func (s *service) PersonalIdentifications(sub string) []api.PersonalIdentificationDataV1 {
	return s.personalIdentificationsMap[sub]
}

func newPersonalIdentificationsResponse(
	baseURL string,
	identifications []api.PersonalIdentificationDataV1,
) api.PersonalIdentificationResponseV1 {
	totalPages := 1
	if len(identifications) == 0 {
		totalPages = 0
	}
	resp := api.PersonalIdentificationResponseV1{
		Data: identifications,
		Links: api.Links{
			Self: fmt.Sprintf("%s/customers/v1/personal/identifications", baseURL),
		},
		Meta: api.Meta{
			TotalPages:   int32(totalPages),
			TotalRecords: int32(len(identifications)),
		},
	}

	return resp
}
