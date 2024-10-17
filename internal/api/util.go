package api

import (
	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-open-insurance/internal/oidc"
)

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
	default:
		return nil
	}
}

func requiredScopes(operationID string) []goidc.Scope {
	switch operationID {
	case "CreateConsentV2", "ConsentV2", "DeleteConsentV2":
		return []goidc.Scope{oidc.ScopeConsents}
	case "PersonalIdentificationsV1", "PersonalQualificationsV1",
		"PersonalComplimentaryInfoV1":
		return []goidc.Scope{oidc.ScopeOpenID, oidc.ScopeConsent}
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
