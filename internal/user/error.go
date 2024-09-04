package user

import (
	"net/http"

	"github.com/luikyv/go-opf/internal/opinerr"
)

var (
	errorUserNotFound = opinerr.New("USER_NOT_FOUND", http.StatusNotFound, "could not find user")
)
