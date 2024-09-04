package consent

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/luikyv/go-opf/internal/opinerr"
	"github.com/luikyv/go-opf/internal/slice"
)

func consentID(nameSpace string) string {
	return fmt.Sprintf("%s:%s", nameSpace, uuid.NewString())
}

func validatePermissions(requestedPermissions []Permission) error {

	if len(requestedPermissions) < 1 {
		return opinerr.New("INVALID_PERMISSION", http.StatusBadRequest,
			"at least one permission must be requested")
	}

	if !slice.ContainsAll(permissions, requestedPermissions...) {
		return opinerr.New("INVALID_PERMISSION", http.StatusBadRequest,
			"invalid permission")
	}

	isPhase2 := slice.ContainsAny(permissionsPhase2, requestedPermissions...)
	isPhase3 := slice.ContainsAny(permissionsPhase3, requestedPermissions...)
	if isPhase2 && isPhase3 {
		return opinerr.New("INVALID_PERMISSION", http.StatusUnprocessableEntity,
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

func validatePermissionsPhase2(requestedPermissions []Permission) error {

	if !slices.Contains(requestedPermissions, PermissionResourcesRead) {
		return opinerr.New("INVALID_PERMISSION", http.StatusBadRequest,
			fmt.Sprintf("the permission %s is required for phase 2", PermissionResourcesRead))
	}

	// RESOURCES_READ cannot be the only permission requested.
	if len(requestedPermissions) == 1 {
		return opinerr.New("INVALID_PERMISSION", http.StatusBadRequest,
			fmt.Sprintf("the permission %s cannot be requested alone", PermissionResourcesRead))
	}

	return nil
}

func validatePermissionsPhase3(requestedPermissions []Permission) error {
	categories := categories(requestedPermissions)

	if len(categories) != 1 {
		return opinerr.New("INVALID_PERMISSION", http.StatusUnprocessableEntity,
			"permissions of different phase 3 categories were requested")
	}

	if !slice.ContainsAll(requestedPermissions, categories[0]...) {
		return opinerr.New("INVALID_PERMISSION", http.StatusBadRequest,
			"all the permission from one category must be requested")
	}

	return nil
}

func categories(requestedPermissions []Permission) []PermissionCategory {
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
