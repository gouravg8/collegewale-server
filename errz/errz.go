package errz

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type BadRequest struct {
	msg string
}

func NewBadRequest(msg string) *BadRequest { return &BadRequest{msg: msg} }
func (e *BadRequest) Error() string        { return e.msg }

type Unauthorized struct {
	msg string
}

func NewUnauthorized(msg string) *Unauthorized { return &Unauthorized{msg: msg} }
func (e *Unauthorized) Error() string          { return e.msg }

type NotFound struct{ msg string }

func NewNotFound(msg string) *NotFound { return &NotFound{msg: msg} }
func (e *NotFound) Error() string      { return e.msg }

type NotAllowed struct{ msg string }

func NewNotAllowed(msg string) *NotAllowed { return &NotAllowed{msg: msg} }
func (e *NotAllowed) Error() string        { return e.msg }

type AlreadyExists struct{ msg string }

func NewAlreadyExists(msg string) *AlreadyExists { return &AlreadyExists{msg: msg} }
func (e *AlreadyExists) Error() string           { return e.msg }

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
