package consent

import (
	"time"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type Consent struct {
	ID              string                  `bson:"_id"`
	Status          api.ConsentStatus       `bson:"status"`
	UserCPF         string                  `bson:"user_cpf"`
	BusinessCNPJ    string                  `bson:"business_cnpj,omitempty"`
	ClientId        string                  `bson:"client_id"`
	Permissions     []api.ConsentPermission `bson:"permissions"`
	CreatedAt       time.Time               `bson:"created_at"`
	UpdatedAt       time.Time               `bson:"updated_at"`
	ExpiresAt       time.Time               `bson:"expires_at"`
	RejectionInfo   *RejectionInfo          `bson:"rejection,omitempty"`
	EndorsementInfo *EndorsementInfo        `bson:"endorsement,omitempty"`
}

// HasAuthExpired returns true if the status is [StatusAwaitingAuthorisation] and
// the max time awaiting authorization has elapsed.
func (c Consent) HasAuthExpired() bool {
	now := time.Now().UTC()
	return c.Status == api.ConsentStatusAWAITINGAUTHORISATION &&
		now.After(c.CreatedAt.Add(time.Second*maxTimeAwaitingAuthorizationSecs))
}

// IsExpired returns true if the status is [StatusAuthorised] and the consent
// reached the expiration date.
func (c Consent) IsExpired() bool {
	now := time.Now().UTC()
	return c.Status == api.ConsentStatusAUTHORISED && now.After(c.ExpiresAt)
}

func (c Consent) IsAuthorized() bool {
	return c.Status == api.ConsentStatusAUTHORISED
}

func (c Consent) HasPermissions(permissions []api.ConsentPermission) bool {
	return containsAll(c.Permissions, permissions...)
}

type RejectionInfo struct {
	RejectedBy api.ConsentRejectedBy         `bson:"rejected_by"`
	Reason     api.ConsentRejectedReasonCode `bson:"reason"`
}

type EndorsementInfo struct {
	PolicyNumber string              `bson:"policy_number"`
	Type         api.EndorsementType `bson:"type"`
	Description  string              `bson:"description"`
}
