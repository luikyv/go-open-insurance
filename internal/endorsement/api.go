package endorsement

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

func (s ServerV1) CreateEndorsementV1(
	ctx context.Context,
	request api.CreateEndorsementV1RequestObject,
) (
	api.CreateEndorsementV1ResponseObject,
	error,
) {
	meta := api.NewRequestMeta(ctx)
	resp, err := s.service.create(ctx, meta, request.ConsentId, *request.Body)
	if err != nil {
		return nil, err
	}

	return api.CreateEndorsementV1201JSONResponse(resp), nil
}
