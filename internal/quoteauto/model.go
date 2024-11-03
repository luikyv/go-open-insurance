package quoteauto

import (
	"time"

	"github.com/google/uuid"
	"github.com/luikyv/go-open-insurance/internal/api"
)

type Lead struct {
	ID                   string                `bson:"_id"`
	ConsentID            string                `bson:"consent_id"`
	Status               api.QuoteStatus       `bson:"status"`
	RejectionReason      string                `bson:"rejection_reason"`
	StatusUpdateDateTime time.Time             `bson:"updated_at"`
	Data                 api.QuoteAutoLeadData `bson:"data"`
}

type Quote struct {
	ID                   string            `bson:"_id"`
	ConsentID            string            `bson:"consent_id"`
	Status               api.QuoteStatus   `bson:"status"`
	RejectionReason      string            `bson:"rejection_reason"`
	StatusUpdateDateTime time.Time         `bson:"updated_at"`
	Data                 api.QuoteAutoData `bson:"data"`
}

func newLead(req api.CreateQuoteAutoLeadRequest) Lead {
	lead := Lead{
		ID:                   uuid.NewString(),
		ConsentID:            req.Data.ConsentId,
		Status:               api.QuoteStatusRCVD,
		StatusUpdateDateTime: time.Now().UTC(),
		Data:                 req.Data,
	}

	return lead
}

func newLeadCreateResponse(
	meta api.RequestMeta,
	lead Lead,
) api.CreateQuoteLeadResponse {
	return api.CreateQuoteLeadResponse{
		Data: api.QuoteStatusInfo{
			Status:               lead.Status,
			StatusUpdateDateTime: api.NewDateTime(lead.StatusUpdateDateTime),
		},
		Links: api.Links{
			Self: meta.RequestURI,
		},
		Meta: api.Meta{
			TotalPages:   1,
			TotalRecords: 1,
		},
	}
}

func newRevokeLeadResponse() api.RevokeQuoteLeadResponse {
	return api.RevokeQuoteLeadResponse{
		Data: struct {
			Status api.RevokeQuoteLeadResponseDataStatus "json:\"status\""
		}{
			Status: api.RevokeQuoteLeadResponseDataStatusCANC,
		},
	}
}

func newQuote(req api.CreateQuoteAutoRequest) Quote {
	quote := Quote{
		ID:                   uuid.NewString(),
		ConsentID:            req.Data.ConsentId,
		Status:               api.QuoteStatusRCVD,
		StatusUpdateDateTime: time.Now().UTC(),
		Data:                 req.Data,
	}

	return quote
}

func newCreateQuoteResponse(
	meta api.RequestMeta,
	quote Quote,
) api.CreateQuoteResponse {
	return api.CreateQuoteResponse{
		Data: api.QuoteStatusInfo{
			Status:               quote.Status,
			StatusUpdateDateTime: api.NewDateTime(quote.StatusUpdateDateTime),
		},
		Links: api.Links{
			Self: meta.Host + "/open-insurance/quote-auto/v1/request/" + quote.ConsentID + "/quote-status",
		},
		Meta: api.Meta{
			TotalPages:   1,
			TotalRecords: 1,
		},
	}
}

func newGetQuoteAutoStatusResponse(
	meta api.RequestMeta,
	quote Quote,
) api.GetQuoteAutoStatusResponse {
	resp := api.GetQuoteAutoStatusResponse{
		Links: api.Links{
			Self: meta.RequestURL(),
		},
		Meta: api.Meta{
			TotalPages:   1,
			TotalRecords: 1,
		},
	}

	resp.Data.Status = api.GetQuoteAutoStatusResponseDataStatus(quote.Status)
	resp.Data.StatusUpdateDateTime = api.NewDateTime(quote.StatusUpdateDateTime)
	if quote.Status != api.QuoteStatusACPT && quote.Status != api.QuoteStatusACKN {
		return resp
	}

	quoteCustomer := quote.Data.QuoteCustomer.ToQuoteCustomer()
	resp.Data.QuoteInfo = &api.QuoteStatusAuto{
		QuoteCustomData: quote.Data.QuoteCustomData,
		QuoteCustomer:   &quoteCustomer,
		Quotes: []struct {
			Assistances         []api.QuoteResultAssistance         "json:\"assistances\""
			Coverages           *[]api.QuoteAutoQuoteResultCoverage "json:\"coverages,omitempty\""
			InsurerQuoteId      string                              "json:\"insurerQuoteId\""
			PremiumInfo         api.QuoteResultPremium              "json:\"premiumInfo\""
			SusepProcessNumbers []string                            "json:\"susepProcessNumbers\""
		}{
			{
				InsurerQuoteId:      quote.ID,
				SusepProcessNumbers: []string{"123456789"},
				Assistances:         []api.QuoteResultAssistance{},
				PremiumInfo: api.QuoteResultPremium{
					PaymentsQuantity: 6,
					TotalPremiumAmount: api.AmountDetails{
						Amount: "100.00",
						Unit: struct {
							Code        string "json:\"code\""
							Description string "json:\"description\""
						}{
							Code:        "BR",
							Description: "BRL",
						},
					},
					TotalNetAmount: api.AmountDetails{
						Amount: "100.00",
						Unit: struct {
							Code        string "json:\"code\""
							Description string "json:\"description\""
						}{
							Code:        "BR",
							Description: "BRL",
						},
					},
					IOF: api.AmountDetails{
						Amount: "100.00",
						Unit: struct {
							Code        string "json:\"code\""
							Description string "json:\"description\""
						}{
							Code:        "BR",
							Description: "BRL",
						},
					},
					Coverages: []api.QuoteResultPremiumCoverage{},
					Payments:  []api.QuoteResultPayment{},
				},
			},
		},
	}

	return resp
}

func newPatchQuoteResponse(
	meta api.RequestMeta,
	quote Quote,
) api.PatchQuoteResponse {
	resp := api.PatchQuoteResponse{}
	resp.Data.Status = api.PatchQuoteResponseDataStatus(quote.Status)
	if quote.Status == api.QuoteStatusCANC {
		return resp
	}

	resp.Data.InsurerQuoteId = &quote.ID
	resp.Data.Links = &api.RedirectLinks{
		// TODO: Change this.
		Redirect: meta.Host + "/auth/.well-known/openid-configuration",
	}
	return resp
}
