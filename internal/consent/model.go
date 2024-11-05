package consent

import (
	"fmt"
	"time"

	"github.com/luikyv/go-open-insurance/internal/api"
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
	RejectionInfo *RejectionInfo          `bson:"rejection,omitempty"`
	Data          api.ConsentData         `json:"data"`
}

// HasAuthExpired returns true if the status is [StatusAwaitingAuthorisation] and
// the max time awaiting authorization has elapsed.
func (c Consent) HasAuthExpired() bool {
	now := time.Now().UTC()
	return c.IsAwaitingAuthorization() &&
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

func (c Consent) IsAwaitingAuthorization() bool {
	return c.Status == api.ConsentStatusAWAITINGAUTHORISATION
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

func newConsent(
	meta api.RequestMeta,
	req api.CreateConsentRequest,
) Consent {
	now := time.Now().UTC()
	c := Consent{
		ID:          ID(),
		Status:      api.ConsentStatusAWAITINGAUTHORISATION,
		UserCPF:     req.Data.LoggedUser.Document.Identification,
		ClientId:    meta.ClientID,
		Permissions: req.Data.Permissions,
		CreatedAt:   now,
		UpdatedAt:   now,
		ExpiresAt:   req.Data.ExpirationDateTime.Time,
		Data:        req.Data,
	}

	if req.Data.BusinessEntity != nil {
		c.BusinessCNPJ = req.Data.BusinessEntity.Document.Identification
	}

	return c
}

func newResponse(meta api.RequestMeta, consent Consent) api.ConsentResponse {
	resp := api.ConsentResponse{
		Meta: &api.Meta{
			TotalRecords: 1,
			TotalPages:   1,
		},
		Links: &api.Links{
			Self: fmt.Sprintf("%s/open-insurance/consents/v2/consents/%s", meta.Host, consent.ID),
		},
	}

	resp.Data.ConsentId = consent.ID
	resp.Data.Status = consent.Status
	resp.Data.Permissions = consent.Permissions
	resp.Data.CreationDateTime = api.NewDateTime(consent.CreatedAt)
	resp.Data.StatusUpdateDateTime = api.NewDateTime(consent.UpdatedAt)
	resp.Data.ExpirationDateTime = api.NewDateTime(consent.ExpiresAt)
	resp.Data.EndorsementInformation = consent.Data.EndorsementInformation

	if consent.RejectionInfo != nil {
		resp.Data.Rejection = &struct {
			Reason     api.ConsentRejectedReason "json:\"reason\""
			RejectedBy api.ConsentRejectedBy     "json:\"rejectedBy\""
		}{
			RejectedBy: consent.RejectionInfo.RejectedBy,
			Reason: api.ConsentRejectedReason{
				Code: consent.RejectionInfo.Reason,
			},
		}
	}

	return resp
}
