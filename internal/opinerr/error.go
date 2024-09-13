package opinerr

import (
	"fmt"
	"net/http"
)

var (
	ErrInternal = Error{"INTERNAL_ERROR", http.StatusInternalServerError, "internal error"}
)

type Error struct {
	Code        string
	StatusCode  int
	Description string
}

func (err Error) Error() string {
	return fmt.Sprintf("%s %s", err.Code, err.Description)
}

func New(code string, status int, description string) Error {
	return Error{
		Code:        code,
		Description: description,
	}
}
