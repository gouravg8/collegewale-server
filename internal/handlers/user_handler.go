package handlers

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	us *service.UserService
}

func NewUserHandler(group *echo.Group, us *service.UserService) *UserHandler {
	h := &UserHandler{us}
	group.GET("/myinfo", h.MyInfo)
	return h
}

func (h UserHandler) MyInfo(ctx echo.Context) error {
	cc := ctx.(*CustomContext)
	if cc == nil {
		return ctx.JSON(http.StatusOK, errz.NewNotFound("user not found"))
	}
	res, err := h.us.MyInfo(cc.user)
	if err != nil {
		return errz.HandleErrx(ctx, err)
	}
	return ctx.JSON(http.StatusOK, res)
}
