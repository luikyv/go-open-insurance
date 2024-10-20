package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/endorsement"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
)

type Server struct {
	service endorsement.Service
}

func NewServer(service endorsement.Service) Server {
	return Server{
		service: service,
	}
}

func (s Server) CreateEndorsementV1(
	ctx context.Context,
	request api.CreateEndorsementV1RequestObject,
) (
	api.CreateEndorsementV1ResponseObject,
	error,
) {
	consentID := ctx.Value(api.CtxKeyConsentID).(string)
	if request.ConsentId != consentID {
		return nil, opinerr.New(
			"NAO_INFORMADO",
			http.StatusBadRequest,
			"invalid consent id",
		)
	}

	sub := ctx.Value(api.CtxKeySubject).(string)
	endorsement := newEndorsement(*request.Body, consentID)
	if err := s.service.Create(ctx, sub, endorsement); err != nil {
		return nil, err
	}

	resp := newResponse(endorsement)
	return api.CreateEndorsementV1201JSONResponse(resp), nil
}

func newEndorsement(
	req api.CreateEndorsementRequestV1,
	consentID string,
) endorsement.Endorsement {
	return endorsement.Endorsement{
		ID:           endorsement.ID(),
		PolicyNumber: req.Data.PolicyNumber,
		ConsentID:    consentID,
		Type:         req.Data.EndorsementType,
		Description:  req.Data.RequestDescription,
		CreatedAt:    time.Now().UTC(),
		RequestedAt:  req.Data.RequestDate.Time,
		CustomData:   req.Data.CustomData,
	}
}

func newResponse(
	endorsement endorsement.Endorsement,
) api.EndorsementResponseV1 {
	return api.EndorsementResponseV1{
		Data: api.EndorsementData{
			ProtocolNumber:     endorsement.ID,
			PolicyNumber:       endorsement.PolicyNumber,
			EndorsementType:    endorsement.Type,
			RequestDescription: endorsement.Description,
			ProtocolDateTime:   api.NewDateTime(endorsement.CreatedAt),
			RequestDate:        api.NewDate(endorsement.RequestedAt),
			CustomData:         endorsement.CustomData,
		},
		Links: api.RedirectLinks{
			Redirect: "https://random.com",
		},
	}
}
