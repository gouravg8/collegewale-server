package handlers

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums/roles"
	"collegeWaleServer/internal/service"
	"collegeWaleServer/internal/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	us *service.UserService
}

func NewUserHandler(group *echo.Group, us *service.UserService) *UserHandler {
	h := &UserHandler{us}
	group.GET("/myinfo", h.MyInfo)
	group.PUT("/user/update", WithRole(h.UpdateUser, roles.Admin))
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

func (h UserHandler) UpdateUser(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc == nil {
		return c.JSON(http.StatusOK, errz.NewNotFound("user not found"))
	}

	var req views.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return errz.HandleErrx(c, err)
	}

	res, err := h.us.UpdateUser(cc.user, req)
	if err != nil {
		return errz.HandleErrx(c, err)
	}

	return c.JSON(http.StatusOK, res)
}
