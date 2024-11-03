package consent

import (
	"fmt"
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
	}

	if req.Data.BusinessEntity != nil {
		c.BusinessCNPJ = req.Data.BusinessEntity.Document.Identification
	}

	if req.Data.EndorsementInformation != nil {
		c.EndorsementInfo = &EndorsementInfo{
			PolicyNumber: req.Data.EndorsementInformation.PolicyNumber,
			Type:         req.Data.EndorsementInformation.EndorsementType,
			Description:  req.Data.EndorsementInformation.RequestDescription,
		}
	}

	return c
}

func newResponse(meta api.RequestMeta, consent Consent) api.ConsentResponse {

	data := api.ConsentDataV2{
		ConsentId:            consent.ID,
		Status:               consent.Status,
		Permissions:          consent.Permissions,
		CreationDateTime:     api.NewDateTime(consent.CreatedAt),
		StatusUpdateDateTime: api.NewDateTime(consent.CreatedAt),
		ExpirationDateTime:   api.NewDateTime(consent.ExpiresAt),
	}

	if consent.RejectionInfo != nil {
		data.Rejection = &api.ConsentRejection{
			RejectedBy: consent.RejectionInfo.RejectedBy,
			Reason: api.ConsentRejectedReason{
				Code: consent.RejectionInfo.Reason,
			},
		}
	}

	if consent.EndorsementInfo != nil {
		data.EndorsementInformation = &api.EndorsementInfo{
			PolicyNumber:       consent.EndorsementInfo.PolicyNumber,
			EndorsementType:    consent.EndorsementInfo.Type,
			RequestDescription: consent.EndorsementInfo.Description,
		}
	}

	return api.ConsentResponse{
		Data: data,
		Meta: &api.Meta{
			TotalRecords: 1,
			TotalPages:   1,
		},
		Links: &api.Links{
			Self: fmt.Sprintf("%s/open-insurance/consents/v2/consents/%s", meta.Host, consent.ID),
		},
	}
}
