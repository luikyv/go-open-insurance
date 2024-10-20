package api

import (
	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-open-insurance/internal/oidc"
)

func isIdempotent(operationID string) bool {
	switch operationID {
	case "CreateEndorsementV1":
		return true
	default:
		return false
	}
}

func isFAPIIDRequired(operationID string) bool {
	switch operationID {
	case "CreateEndorsementV1":
		return true
	default:
		return false
	}
}

func requiredPermissions(operationID string) []ConsentPermission {
	switch operationID {
	case "PersonalIdentificationsV1":
		return []ConsentPermission{
			ConsentPermissionRESOURCESREAD,
			ConsentPermissionCUSTOMERSPERSONALIDENTIFICATIONSREAD,
		}
	case "PersonalQualificationsV1":
		return []ConsentPermission{
			ConsentPermissionRESOURCESREAD,
			ConsentPermissionCUSTOMERSPERSONALQUALIFICATIONREAD,
		}
	case "PersonalComplimentaryInfoV1":
		return []ConsentPermission{
			ConsentPermissionRESOURCESREAD,
			ConsentPermissionCUSTOMERSPERSONALADDITIONALINFOREAD,
		}
	case "ResourcesV2":
		return []ConsentPermission{ConsentPermissionRESOURCESREAD}
	case "CapitalizationTitlePlansV1":
		return []ConsentPermission{
			ConsentPermissionRESOURCESREAD,
			ConsentPermissionCAPITALIZATIONTITLEREAD,
		}
	case "CapitalizationTitlePlanInfoV1":
		return []ConsentPermission{
			ConsentPermissionRESOURCESREAD,
			ConsentPermissionCAPITALIZATIONTITLEPLANINFOREAD,
		}
	case "CapitalizationTitleEventsV1":
		return []ConsentPermission{
			ConsentPermissionRESOURCESREAD,
			ConsentPermissionCAPITALIZATIONTITLEEVENTSREAD,
		}
	case "CapitalizationTitleSettlementsV1":
		return []ConsentPermission{
			ConsentPermissionRESOURCESREAD,
			ConsentPermissionCAPITALIZATIONTITLESETTLEMENTSREAD,
		}
	case "CreateEndorsementV1":
		return []ConsentPermission{
			ConsentPermissionENDORSEMENTREQUESTCREATE,
		}
	default:
		return nil
	}
}

func requiredScopes(operationID string) []goidc.Scope {
	switch operationID {
	case "CreateConsentV2", "ConsentV2", "DeleteConsentV2":
		return []goidc.Scope{oidc.ScopeConsents}
	case "ResourcesV2":
		return []goidc.Scope{oidc.ScopeOpenID, oidc.ScopeConsent, oidc.ScopeResources}
	case "PersonalIdentificationsV1", "PersonalQualificationsV1",
		"PersonalComplimentaryInfoV1":
		return []goidc.Scope{oidc.ScopeOpenID, oidc.ScopeConsent, oidc.ScopeCustomers}
	case "CapitalizationTitlePlansV1", "CapitalizationTitlePlanInfoV1",
		"CapitalizationTitleEventsV1", "CapitalizationTitleSettlementsV1":
		return []goidc.Scope{oidc.ScopeOpenID, oidc.ScopeConsent, oidc.ScopeCapitalizationTitle}
	case "CreateEndorsementV1":
		return []goidc.Scope{oidc.ScopeOpenID, oidc.ScopeConsent, oidc.ScopeEndorsement}
	default:
		return nil
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
