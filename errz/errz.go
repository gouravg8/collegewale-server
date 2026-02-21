package errz

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type BadRequest struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewBadRequest(msg string) *BadRequest { return &BadRequest{msg, http.StatusBadRequest} }
func (e *BadRequest) Error() string        { return e.Message }

type Unauthorized struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewUnauthorized(msg string) *Unauthorized { return &Unauthorized{msg, http.StatusUnauthorized} }
func (e *Unauthorized) Error() string          { return e.Message }

type NotFound struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewNotFound(msg string) *NotFound { return &NotFound{msg, http.StatusNotFound} }
func (e *NotFound) Error() string      { return e.Message }

type Forbiddenn struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewForbidden(msg string) *Forbiddenn { return &Forbiddenn{msg, http.StatusForbidden} }
func (e *Forbiddenn) Error() string       { return e.Message }

type AlreadyExists struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewAlreadyExists(msg string) *AlreadyExists { return &AlreadyExists{msg, http.StatusConflict} }
func (e *AlreadyExists) Error() string           { return e.Message }

func HandleErrz[T any](ctx echo.Context, res T, err error) error {
	if err == nil {
		return ctx.JSON(http.StatusOK, res)
	}
	var existsErr *AlreadyExists
	var notFoundErr *NotFound
	var authErr *Unauthorized
	var badReqErr *BadRequest
	var notAllowed *Forbiddenn
	switch {
	case errors.As(err, &notAllowed):
		return ctx.JSON(http.StatusForbidden, err)

	case errors.As(err, &existsErr):
		return ctx.JSON(http.StatusConflict, err)

	case errors.As(err, &notFoundErr):
		return ctx.JSON(http.StatusNotFound, err)

	case errors.As(err, &authErr):
		return ctx.JSON(http.StatusUnauthorized, err)

	case errors.As(err, &badReqErr):
		return ctx.JSON(http.StatusBadRequest, err)

	default:
		log.Errorf("something went wrong :: %+v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"message": fmt.Sprintf("something went wrong :: %v", err), "code": http.StatusInternalServerError})
	}
}

func HandleErrx(ctx echo.Context, err error) error {
	if err == nil {
		return ctx.JSON(http.StatusOK, nil)
	}
	var existsErr *AlreadyExists
	var notFoundErr *NotFound
	var authErr *Unauthorized
	var badReqErr *BadRequest
	var notAllowed *Forbiddenn
	switch {
	case errors.As(err, &notAllowed):
		return ctx.JSON(http.StatusForbidden, err)

	case errors.As(err, &existsErr):
		return ctx.JSON(http.StatusConflict, err)

	case errors.As(err, &notFoundErr):
		return ctx.JSON(http.StatusNotFound, err)

	case errors.As(err, &authErr):
		return ctx.JSON(http.StatusUnauthorized, err)

	case errors.As(err, &badReqErr):
		return ctx.JSON(http.StatusBadRequest, err)

	default:
		log.Errorf("something went wrong :: %+v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"message": fmt.Sprintf("something went wrong :: %v", err), "code": http.StatusInternalServerError})
	}
}
