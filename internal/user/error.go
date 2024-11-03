package user

import (
	"net/http"

	"github.com/luikyv/go-open-insurance/internal/api"
)

var (
	errorUserNotFound = api.NewError("USER_NOT_FOUND", http.StatusNotFound, "could not find user")
)
