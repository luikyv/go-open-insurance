package customer

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

func (s ServerV1) PersonalIdentificationsV1(
	ctx context.Context,
	request api.PersonalIdentificationsV1RequestObject,
) (
	api.PersonalIdentificationsV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp := s.service.personalIdentifications(meta)
	return api.PersonalIdentificationsV1200JSONResponse(resp), nil
}

func (s ServerV1) PersonalQualificationsV1(
	ctx context.Context,
	request api.PersonalQualificationsV1RequestObject,
) (
	api.PersonalQualificationsV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp := s.service.personalQualifications(meta)
	return api.PersonalQualificationsV1200JSONResponse(resp), nil
}

func (s ServerV1) PersonalComplimentaryInfoV1(
	ctx context.Context,
	request api.PersonalComplimentaryInfoV1RequestObject,
) (
	api.PersonalComplimentaryInfoV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp := s.service.personalComplimentaryInfos(meta)
	return api.PersonalComplimentaryInfoV1200JSONResponse(resp), nil
}
