package api

import (
	"github.com/luikyv/go-oidc/pkg/goidc"
)

type operationOptions struct {
	scopes           []goidc.Scope
	permissions      []ConsentPermission
	fapiIDIsRequired bool
	isIdempotent     bool
}

func newOperationOptions(operationID string) operationOptions {
	switch operationID {
	case "CreateConsentV2", "ConsentV2", "DeleteConsentV2":
		return operationOptions{
			scopes: []goidc.Scope{ScopeConsents},
		}
	case "ResourcesV2":
		return operationOptions{
			scopes: []goidc.Scope{
				ScopeOpenID,
				ScopeConsent,
				ScopeResources,
			},
			permissions: []ConsentPermission{ConsentPermissionRESOURCESREAD},
		}
	case "PersonalIdentificationsV1":
		return operationOptions{
			scopes: []goidc.Scope{
				ScopeOpenID,
				ScopeConsent,
				ScopeCustomers,
			},
			permissions: []ConsentPermission{
				ConsentPermissionRESOURCESREAD,
				ConsentPermissionCUSTOMERSPERSONALIDENTIFICATIONSREAD,
			},
		}
	case "PersonalQualificationsV1":
		return operationOptions{
			scopes: []goidc.Scope{
				ScopeOpenID,
				ScopeConsent,
				ScopeCustomers,
			},
			permissions: []ConsentPermission{
				ConsentPermissionRESOURCESREAD,
				ConsentPermissionCUSTOMERSPERSONALQUALIFICATIONREAD,
			},
		}
	case "PersonalComplimentaryInfoV1":
		return operationOptions{
			scopes: []goidc.Scope{
				ScopeOpenID,
				ScopeConsent,
				ScopeCustomers,
			},
			permissions: []ConsentPermission{
				ConsentPermissionRESOURCESREAD,
				ConsentPermissionCUSTOMERSPERSONALADDITIONALINFOREAD,
			},
		}
	case "CapitalizationTitlePlansV1":
		return operationOptions{
			scopes: []goidc.Scope{
				ScopeOpenID,
				ScopeConsent,
				ScopeCapitalizationTitle,
			},
			permissions: []ConsentPermission{
				ConsentPermissionRESOURCESREAD,
				ConsentPermissionCAPITALIZATIONTITLEREAD,
			},
		}
	case "CapitalizationTitlePlanInfoV1":
		return operationOptions{
			scopes: []goidc.Scope{
				ScopeOpenID,
				ScopeConsent,
				ScopeCapitalizationTitle,
			},
			permissions: []ConsentPermission{
				ConsentPermissionRESOURCESREAD,
				ConsentPermissionCAPITALIZATIONTITLEPLANINFOREAD,
			},
		}
	case "CapitalizationTitleEventsV1":
		return operationOptions{
			scopes: []goidc.Scope{
				ScopeOpenID,
				ScopeConsent,
				ScopeCapitalizationTitle,
			},
			permissions: []ConsentPermission{
				ConsentPermissionRESOURCESREAD,
				ConsentPermissionCAPITALIZATIONTITLEEVENTSREAD,
			},
		}
	case "CapitalizationTitleSettlementsV1":
		return operationOptions{
			scopes: []goidc.Scope{
				ScopeOpenID,
				ScopeConsent,
				ScopeCapitalizationTitle,
			},
			permissions: []ConsentPermission{
				ConsentPermissionRESOURCESREAD,
				ConsentPermissionCAPITALIZATIONTITLESETTLEMENTSREAD,
			},
		}
	case "CreateEndorsementV1":
		return operationOptions{
			scopes: []goidc.Scope{
				ScopeOpenID,
				ScopeConsent,
				ScopeEndorsement,
			},
			permissions: []ConsentPermission{
				ConsentPermissionENDORSEMENTREQUESTCREATE,
			},
			fapiIDIsRequired: true,
			isIdempotent:     true,
		}
	default:
		return operationOptions{}
	}
}

// areScopesValid verifies every scope in requiredScopes has a match among
// scopes.
// scopes can have more scopes than the defined at requiredScopes, but the
// contrary results in false.
func areScopesValid(requiredScopes []goidc.Scope, scopes []string) bool {
	for _, requiredScope := range requiredScopes {
		if !isScopeValid(requiredScope, scopes) {
			return false
		}
	}
	return true
}

// isScopeValid verifies if requireScope has a match in scopes.
func isScopeValid(requiredScope goidc.Scope, scopes []string) bool {
	for _, scope := range scopes {
		if requiredScope.Matches(scope) {
			return true
		}
	}

	return false
}
