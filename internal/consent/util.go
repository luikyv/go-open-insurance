package consent

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/luikyv/go-open-insurance/internal/api"
)

func ID() string {
	return fmt.Sprintf("urn:mockin:%s", uuid.NewString())
}

func validate(ctx context.Context, consent Consent) error {
	now := time.Now().UTC()
	if now.After(consent.ExpiresAt) {
		return api.NewError("INVALID_REQUEST", http.StatusBadRequest,
			"the expiration time cannot be in the past")
	}

	if consent.ExpiresAt.After(now.AddDate(1, 0, 0)) {
		return api.NewError("INVALID_REQUEST", http.StatusBadRequest,
			"the expiration time cannot be greater than one year")
	}

	if err := validatePermissions(ctx, consent.Permissions); err != nil {
		return err
	}

	if slices.Contains(consent.Permissions, api.ConsentPermissionENDORSEMENTREQUESTCREATE) &&
		consent.Data.EndorsementInformation == nil {
		return api.NewError("INVALID_REQUEST", http.StatusBadRequest,
			"endorsement information is missing")
	}

	return nil
}

func validatePermissions(_ context.Context, requestedPermissions []api.ConsentPermission) error {

	isPhase2 := containsAny(permissionsPhase2, requestedPermissions...)
	isPhase3 := containsAny(permissionsPhase3, requestedPermissions...)
	if isPhase2 && isPhase3 {
		return api.NewError("NAO_INFORMADO", http.StatusUnprocessableEntity,
			"cannot request permission from phase 2 and 3 in the same request")
	}

	if isPhase2 {
		return validatePermissionsPhase2(requestedPermissions)
	}

	if isPhase3 {
		return validatePermissionsPhase3(requestedPermissions)
	}

	return nil
}

func validatePermissionsPhase2(requestedPermissions []api.ConsentPermission) error {

	if !slices.Contains(requestedPermissions, api.ConsentPermissionRESOURCESREAD) {
		return api.NewError("NAO_INFORMADO", http.StatusBadRequest,
			fmt.Sprintf("the permission %s is required for phase 2", api.ConsentPermissionRESOURCESREAD))
	}

	// RESOURCES_READ cannot be the only permission requested.
	if len(requestedPermissions) == 1 {
		return api.NewError("NAO_INFORMADO", http.StatusBadRequest,
			fmt.Sprintf("the permission %s cannot be requested alone", api.ConsentPermissionRESOURCESREAD))
	}

	return nil
}

func validatePermissionsPhase3(requestedPermissions []api.ConsentPermission) error {
	categories := categories(requestedPermissions)

	if len(categories) != 1 {
		return api.NewError("NAO_INFORMADO", http.StatusUnprocessableEntity,
			"permissions of different phase 3 categories were requested")
	}

	if !containsAll(requestedPermissions, categories[0]...) {
		return api.NewError("NAO_INFORMADO", http.StatusBadRequest,
			"all the permission from one category must be requested")
	}

	return nil
}

func categories(requestedPermissions []api.ConsentPermission) []PermissionCategory {
	var categories []PermissionCategory
	for _, cat := range permissionCategories {
		for _, p := range requestedPermissions {
			if cat.contains(p) {
				categories = append(categories, cat)
			}
		}
	}

	return categories
}

func containsAll[T comparable](superSet []T, subSet ...T) bool {
	for _, t := range subSet {
		if !slices.Contains(superSet, t) {
			return false
		}
	}

	return true
}

func containsAny[T comparable](superSet []T, subSet ...T) bool {
	for _, t := range subSet {
		if slices.Contains(superSet, t) {
			return true
		}
	}

	return false
}
