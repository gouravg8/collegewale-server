package errz

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type BadRequest struct {
	Message string `json:"message"`
}

func NewBadRequest(msg string) *BadRequest { return &BadRequest{msg} }
func (e *BadRequest) Error() string        { return e.Message }

type Unauthorized struct {
	Message string `json:"message"`
}

func NewUnauthorized(msg string) *Unauthorized { return &Unauthorized{msg} }
func (e *Unauthorized) Error() string          { return e.Message }

type NotFound struct {
	Message string `json:"message"`
}

func NewNotFound(msg string) *NotFound { return &NotFound{msg} }
func (e *NotFound) Error() string      { return e.Message }

type NotAllowed struct {
	Message string `json:"message"`
}

func NewNotAllowed(msg string) *NotAllowed { return &NotAllowed{msg} }
func (e *NotAllowed) Error() string        { return e.Message }

type AlreadyExists struct {
	Message string `json:"message"`
}

func NewAlreadyExists(msg string) *AlreadyExists { return &AlreadyExists{msg} }
func (e *AlreadyExists) Error() string           { return e.Message }

func HandleErrz[T any](ctx echo.Context, res T, err error) error {
	if err == nil {
		return ctx.JSON(http.StatusOK, res)
	}
	var existsErr *AlreadyExists
	var notFoundErr *NotFound
	var authErr *Unauthorized
	var badReqErr *BadRequest
	var notAllowed *NotAllowed
	switch {
	case errors.As(err, &notAllowed):
		return ctx.JSON(http.StatusForbidden, err.Error())

	case errors.As(err, &existsErr):
		return ctx.JSON(http.StatusConflict, err.Error())

	case errors.As(err, &notFoundErr):
		return ctx.JSON(http.StatusNotFound, err.Error())

	case errors.As(err, &authErr):
		return ctx.JSON(http.StatusUnauthorized, err.Error())

	case errors.As(err, &badReqErr):
		return ctx.JSON(http.StatusBadRequest, err.Error())

	default:
		log.Errorf("something went wrong :: %+v", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
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
	var notAllowed *NotAllowed
	switch {
	case errors.As(err, &notAllowed):
		return ctx.JSON(http.StatusForbidden, err.Error())

	case errors.As(err, &existsErr):
		return ctx.JSON(http.StatusConflict, err.Error())

	case errors.As(err, &notFoundErr):
		return ctx.JSON(http.StatusNotFound, err.Error())

	case errors.As(err, &authErr):
		return ctx.JSON(http.StatusUnauthorized, err.Error())

	case errors.As(err, &badReqErr):
		return ctx.JSON(http.StatusBadRequest, err.Error())

	default:
		log.Errorf("something went wrong :: %+v", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
}
