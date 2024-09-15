package consent

import (
	"net/http"

	"github.com/luikyv/go-open-insurance/internal/opinerr"
)

var (
	errNotFound = opinerr.New("NOT_FOUND", http.StatusNotFound,
		"could not find the consent")
	errAlreadyRejected = opinerr.New("INVALID_OPERATION", http.StatusBadRequest,
		"the consent is already rejected")
	errInvalidStatus = opinerr.New("INVALID_STATUS", http.StatusBadRequest,
		"invalid consent status")
	errInvalidPermissions = opinerr.New("INVALID_PERMISSIONS", http.StatusBadRequest,
		"invalid consent permissions")
	errClientNotAuthorized = opinerr.New("UNAUTHORIZED", http.StatusForbidden,
		"client not authorized to perform this operation")
)
