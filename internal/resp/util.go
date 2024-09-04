package resp

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luikyv/go-opf/internal/opinerr"
	"github.com/luikyv/go-opf/internal/time"
)

type Links struct {
	Self string `json:"self"`
}

type Meta struct {
	TotalRecords int `json:"totalRecords"`
	TotalPages   int `json:"totalPages"`
}

type Response struct {
	Data  any   `json:"data"`
	Meta  Meta  `json:"meta"`
	Links Links `json:"links"`
}

type Error struct {
	Code            string        `json:"code"`
	Title           string        `json:"title"`
	Detail          string        `json:"detail"`
	RequestDateTime time.DateTime `json:"requestDateTime"`
}

type ResponseError struct {
	Errors []Error `json:"errors"`
	Meta   Meta    `json:"meta"`
}

func newError(code string, description string) ResponseError {
	return ResponseError{
		Errors: []Error{{Code: code, Title: description, Detail: description}},
		Meta: Meta{
			TotalRecords: 1,
			TotalPages:   1,
		},
	}
}

func WriteError(ctx *gin.Context, err error) {
	var opfErr opinerr.Error
	if !errors.As(err, &opfErr) {
		ctx.JSON(
			http.StatusInternalServerError,
			newError(opinerr.ErrorInternal.Code, opinerr.ErrorInternal.Description),
		)
		return
	}

	ctx.JSON(
		opfErr.StatusCode,
		newError(opfErr.Code, opfErr.Description),
	)
}
