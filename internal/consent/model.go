package consent

import (
	"time"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type Status string

const (
	StatusAuthorised            Status = "AUTHORISED"
	StatusAwaitingAuthorisation Status = "AWAITING_AUTHORISATION"
	StatusRejected              Status = "REJECTED"
	StatusConsumed              Status = "CONSUMED"
)

type RejectionReason string

const (
	RejectionReasonConsentExpired                RejectionReason = "CONSENT_EXPIRED"
	RejectionReasonCustomerManuallyRejected      RejectionReason = "CUSTOMER_MANUALLY_REJECTED"
	RejectionReasonCustomerManuallyRevoked       RejectionReason = "CUSTOMER_MANUALLY_REVOKED"
	RejectionReasonConsentMaxDateReached         RejectionReason = "CONSENT_MAX_DATE_REACHED"
	RejectionReasonConsentTechnicalIssue         RejectionReason = "CONSENT_TECHNICAL_ISSUE"
	RejectionReasonConsentInternalSecurityReason RejectionReason = "INTERNAL_SECURITY_REASON"
)

type RejectedBy string

const (
	RejectedByUser  RejectedBy = "USER"
	RejectedByASPSP RejectedBy = "ASPSP"
	RejectedByTPP   RejectedBy = "TPP"
)

type Consent struct {
	ID            string                  `bson:"_id"`
	Status        api.ConsentStatus       `bson:"status"`
	UserCPF       string                  `bson:"user_cpf"`
	BusinessCNPJ  string                  `bson:"business_cnpj,omitempty"`
	ClientId      string                  `bson:"client_id"`
	Permissions   []api.ConsentPermission `bson:"permissions"`
	CreatedAt     time.Time               `bson:"created_at"`
	UpdatedAt     time.Time               `bson:"updated_at"`
	ExpiresAt     time.Time               `bson:"expires_at"`
	RejectionInfo *RejectionInfo          `json:"rejection,omitempty"`
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
	RejectedBy api.ConsentRejectedBy         `json:"rejected_by"`
	Reason     api.ConsentRejectedReasonCode `json:"reason"`
}
