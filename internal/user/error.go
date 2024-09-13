package user

import (
	"net/http"

	"github.com/luikyv/go-open-insurance/internal/opinerr"
)

var (
	errorUserNotFound = opinerr.New("USER_NOT_FOUND", http.StatusNotFound, "could not find user")
)
