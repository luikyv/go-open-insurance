package endorsement

import (
	"time"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type Endorsement struct {
	// ID is the endorsement protocol number.
	ID           string
	PolicyNumber string
	ConsentID    string
	Type         api.EndorsementType
	Description  string
	CreatedAt    time.Time
	RequestedAt  time.Time
	CustomData   *api.EndorsementCustomData
}

func newEndorsement(
	req api.CreateEndorsementRequest,
	consentID string,
) Endorsement {
	return Endorsement{
		ID:           ID(),
		PolicyNumber: req.Data.PolicyNumber,
		ConsentID:    consentID,
		Type:         req.Data.EndorsementType,
		Description:  req.Data.RequestDescription,
		CreatedAt:    time.Now().UTC(),
		RequestedAt:  req.Data.RequestDate.Time,
		CustomData:   req.Data.CustomData,
	}
}

func newCreateResponse(
	endorsement Endorsement,
) api.CreateEndorsementResponse {
	return api.CreateEndorsementResponse{
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
