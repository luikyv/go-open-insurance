package consent

import (
	"fmt"
	"net/http"

	"github.com/luikyv/go-opf/internal/opinerr"
	"github.com/luikyv/go-opf/internal/resp"
	"github.com/luikyv/go-opf/internal/time"
	"github.com/luikyv/go-opf/internal/user"
)

const (
	StatusAuthorised            Status = "AUTHORISED"
	StatusAwaitingAuthorisation Status = "AWAITING_AUTHORISATION"
	StatusRejected              Status = "REJECTED"
	StatusConsumed              Status = "CONSUMED"
)

type Status string

type Consent struct {
	ID                 string        `bson:"_id"`
	Status             Status        `bson:"status"`
	UserCPF            string        `bson:"user_cpf"`
	BusinessCNPJ       string        `bson:"business_cnpj,omitempty"`
	ClientId           string        `bson:"client_id"`
	Permissions        []Permission  `bson:"permissions"`
	CreatedAtTimestamp time.DateTime `bson:"created_at"`
	UpdatedAtTimestamp time.DateTime `bson:"updated_at"`
	ExpiresAtTimestamp time.DateTime `bson:"expires_at"`
}

func (c Consent) IsExpired() bool {
	now := time.Now()

	if c.Status == StatusAwaitingAuthorisation &&
		now.After(c.CreatedAtTimestamp.Add(maxTimeAwaitingAuthorizationSecs)) {
		return true
	}

	if c.Status == StatusAuthorised && now.After(c.ExpiresAtTimestamp) {
		return true
	}

	return false
}

func (c Consent) newResponseV2(baseURL string) resp.Response {
	data := responseDataV2{
		responseData: c.newResponseData(),
	}

	return resp.Response{
		Data: data,
		Meta: resp.Meta{
			TotalRecords: 1,
			TotalPages:   1,
		},
		Links: resp.Links{
			Self: fmt.Sprintf("%s%s/consents/%s", baseURL, apiPrefixConsentsV2, c.ID),
		},
	}
}

func (c Consent) newResponseData() responseData {
	return responseData{
		ConsentID:            c.ID,
		CreationDateTime:     c.CreatedAtTimestamp,
		Status:               c.Status,
		StatusUpdateDateTime: c.CreatedAtTimestamp,
		Permissions:          c.Permissions,
		ExpirationDateTime:   c.ExpiresAtTimestamp,
	}
}

type requestDataV2 struct {
	requestData
}

func newV2(req requestDataV2, clientID, nameSpace string) Consent {
	now := time.Now()
	consent := Consent{
		ID:                 consentID(nameSpace),
		Status:             StatusAwaitingAuthorisation,
		UserCPF:            req.LoggedUser.Document.Identification,
		ClientId:           clientID,
		Permissions:        req.Permissions,
		CreatedAtTimestamp: now,
		UpdatedAtTimestamp: now,
		ExpiresAtTimestamp: req.ExpirationDateTime,
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
	ExpirationDateTime time.DateTime        `json:"expirationDateTime" binding:"required"`
}

func (r requestData) validate() error {
	now := time.Now()
	if now.After(r.ExpirationDateTime) {
		return opinerr.New("INVALID_REQUEST", http.StatusBadRequest,
			"the expiration time cannot be in the past")
	}

	if r.ExpirationDateTime.After(now.AddYears(1)) {
		return opinerr.New("INVALID_REQUEST", http.StatusBadRequest,
			"the expiration time cannot be greater than one year")
	}

	if err := validatePermissions(r.Permissions); err != nil {
		return err
	}

	return nil
}

type responseDataV2 struct {
	responseData
}

type responseData struct {
	ConsentID            string        `json:"consentId"`
	Status               Status        `json:"status"`
	Permissions          []Permission  `json:"permissions"`
	CreationDateTime     time.DateTime `json:"creationDateTime"`
	StatusUpdateDateTime time.DateTime `json:"statusUpdateDateTime"`
	ExpirationDateTime   time.DateTime `json:"expirationDateTime"`
}
