package opinresp

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
	"github.com/luikyv/go-open-insurance/internal/timeutil"
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
	Code            string            `json:"code"`
	Title           string            `json:"title"`
	Detail          string            `json:"detail"`
	RequestDateTime timeutil.DateTime `json:"requestDateTime"`
}

type ResponseError struct {
	Errors []Error `json:"errors"`
	Meta   Meta    `json:"meta"`
}

func newError(err opinerr.Error) ResponseError {
	return ResponseError{
		Errors: []Error{{Code: err.Code, Title: err.Description, Detail: err.Description}},
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
			opinerr.ErrInternal.StatusCode,
			newError(opinerr.ErrInternal),
		)
		return
	}

	ctx.JSON(
		opfErr.StatusCode,
		newError(opfErr),
	)
}
