package consent

import (
	"fmt"

	"github.com/luikyv/go-open-insurance/internal/opinresp"
	"github.com/luikyv/go-open-insurance/internal/timeutil"
	"github.com/luikyv/go-open-insurance/internal/user"
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
	ID            string            `bson:"_id"`
	Status        Status            `bson:"status"`
	UserCPF       string            `bson:"user_cpf"`
	BusinessCNPJ  string            `bson:"business_cnpj,omitempty"`
	ClientId      string            `bson:"client_id"`
	Permissions   []Permission      `bson:"permissions"`
	CreatedAt     timeutil.DateTime `bson:"created_at"`
	UpdatedAt     timeutil.DateTime `bson:"updated_at"`
	ExpiresAt     timeutil.DateTime `bson:"expires_at"`
	RejectionInfo *RejectionInfo    `json:"rejection,omitempty"`
}

// HasAuthExpired returns true if the status is [StatusAwaitingAuthorisation] and
// the max time awaiting authorization has elapsed.
func (c Consent) HasAuthExpired() bool {
	now := timeutil.Now()
	return c.Status == StatusAwaitingAuthorisation &&
		now.After(c.CreatedAt.Add(maxTimeAwaitingAuthorizationSecs))
}

// IsExpired returns true if the status is [StatusAuthorised] and the consent
// reached the expiration date.
func (c Consent) IsExpired() bool {
	now := timeutil.Now()
	return c.Status == StatusAuthorised && now.After(c.ExpiresAt)
}

func (c Consent) response(baseURL string) opinresp.Response {
	var rejection *rejectionResponseData
	if c.RejectionInfo != nil {
		rejection = &rejectionResponseData{}
		rejection.RejectedBy = c.RejectionInfo.RejectedBy
		rejection.Reason.Code = c.RejectionInfo.Reason
	}

	return opinresp.Response{
		Data: responseData{
			ConsentID:            c.ID,
			CreationDateTime:     c.CreatedAt,
			Status:               c.Status,
			StatusUpdateDateTime: c.CreatedAt,
			Permissions:          c.Permissions,
			ExpirationDateTime:   c.ExpiresAt,
			Rejection:            rejection,
		},
		Meta: opinresp.Meta{
			TotalRecords: 1,
			TotalPages:   1,
		},
		Links: opinresp.Links{
			Self: fmt.Sprintf("%s/consents/%s", baseURL, c.ID),
		},
	}
}

func newConsent(req requestData, clientID, nameSpace string) Consent {
	now := timeutil.Now()
	consent := Consent{
		ID:          consentID(nameSpace),
		Status:      StatusAwaitingAuthorisation,
		UserCPF:     req.LoggedUser.Document.Identification,
		ClientId:    clientID,
		Permissions: req.Permissions,
		CreatedAt:   now,
		UpdatedAt:   now,
		ExpiresAt:   req.ExpirationDateTime,
	}

	if req.BusinessEntity != nil {
		consent.BusinessCNPJ = req.BusinessEntity.Document.Identification
	}

	return consent
}

type requestData struct {
	LoggedUser         user.Logged          `json:"loggedUser" binding:"required"`
	BusinessEntity     *user.BusinessEntity `json:"businessEntity,omitempty"`
	Permissions        []Permission         `json:"permissions" binding:"required"`
	ExpirationDateTime timeutil.DateTime    `json:"expirationDateTime" binding:"required"`
}

type responseData struct {
	ConsentID            string                 `json:"consentId"`
	Status               Status                 `json:"status"`
	Permissions          []Permission           `json:"permissions"`
	CreationDateTime     timeutil.DateTime      `json:"creationDateTime"`
	StatusUpdateDateTime timeutil.DateTime      `json:"statusUpdateDateTime"`
	ExpirationDateTime   timeutil.DateTime      `json:"expirationDateTime"`
	Rejection            *rejectionResponseData `json:"rejection,omitempty"`
}

type rejectionResponseData struct {
	RejectedBy RejectedBy `json:"rejectedBy"`
	Reason     struct {
		Code RejectionReason `json:"code"`
	} `json:"reason"`
}

type RejectionInfo struct {
	RejectedBy RejectedBy      `json:"rejected_by"`
	Reason     RejectionReason `json:"reason"`
}
