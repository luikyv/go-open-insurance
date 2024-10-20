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
