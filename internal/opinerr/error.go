package opinerr

import (
	"fmt"
	"net/http"
)

var (
	ErrInternal = Error{"INTERNAL_ERROR", http.StatusInternalServerError, "internal error", false}
)

type Error struct {
	Code           string
	StatusCode     int
	Description    string
	RenderAsSingle bool
}

func (err Error) Error() string {
	return fmt.Sprintf("%s %s", err.Code, err.Description)
}

func New(code string, status int, description string) Error {
	err := Error{
		Code:        code,
		StatusCode:  status,
		Description: description,
	}

	if status == http.StatusUnprocessableEntity {
		err.RenderAsSingle = true
	}

	return err
}
