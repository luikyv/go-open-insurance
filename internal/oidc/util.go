package oidc

import (
	"strings"

	"github.com/luikyv/go-oidc/pkg/goidc"
)

const (
	// ACRs.
	ACROpenInsuranceLOA2 goidc.ACR = "urn:brasil:openbanking:loa2"
	ACROpenInsuranceLOA3 goidc.ACR = "urn:brasil:openbanking:loa3"
)

var (
	ScopeOpenID  = goidc.ScopeOpenID
	ScopeConsent = goidc.NewDynamicScope("consent", func(requestedScope string) bool {
		return strings.HasPrefix(requestedScope, "consent:")
	})
	ScopeConsents                    = goidc.NewScope("consents")
	ScopeResources                   = goidc.NewScope("resources")
	ScopeCustomers                   = goidc.NewScope("customers")
	ScopeAcceptanceAndBranchesAbroad = goidc.NewScope("insurance-acceptance-and-branches-abroad")
	ScopeInsuranceAuto               = goidc.NewScope("insurance-auto")
	ScopeInsuranceFinancialRisk      = goidc.NewScope("insurance-financial-risk")
	ScopeInsurancePatrimonial        = goidc.NewScope("insurance-patrimonial")
	ScopeInsuranceResponsibility     = goidc.NewScope("insurance-responsibility")
)

var Scopes = []goidc.Scope{
	ScopeOpenID,
	ScopeConsent,
	ScopeConsents,
	ScopeResources,
	ScopeCustomers,
	ScopeAcceptanceAndBranchesAbroad,
	ScopeInsuranceAuto,
	ScopeInsuranceFinancialRisk,
	ScopeInsurancePatrimonial,
	ScopeInsuranceResponsibility,
}
