package opinerr

import (
	"fmt"
	"net/http"
)

// TODO: Allow passing multiple errors as a chain.

var (
	ErrorInternal = New("INTERNAL_ERROR", http.StatusInternalServerError,
		"internal error")
)

type Error struct {
	Code        string
	StatusCode  int
	Description string
}

func (err Error) Error() string {
	return fmt.Sprintf("%s %s", err.Code, err.Description)
}

func New(code string, statusCode int, description string) Error {
	return Error{
		Code:        code,
		StatusCode:  statusCode,
		Description: description,
	}
}
