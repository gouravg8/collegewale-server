package handlers

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums/roles"
	service "collegeWaleServer/internal/service/auth"
	"collegeWaleServer/internal/views"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Registry struct {
	s *service.RegistryService
}

func NewRegistryHandler(group *echo.Group, registryService *service.RegistryService) *Registry {
	h := &Registry{
		s: registryService,
	}
	group.POST("/register/college", WithRole(h.RegisterCollege, roles.Admin))
	group.POST("/register/student", WithRole(h.RegisterStudent, roles.Admin, roles.College))
	return h
}

func (h Registry) RegisterCollege(ctx echo.Context) error {
	var req views.College
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errz.NewBadRequest("invalid request"))
	}
	if err := req.IsValidRequest(); err != nil {
		return err
	}
	cc := ctx.(*CustomContext)
	if err := h.s.RegisterCollege(req, cc.user); err != nil {
		return errz.HandleErrx(ctx, err)
	}
	return ctx.JSON(http.StatusOK, views.Response{Message: "success"})
}

func (h Registry) RegisterStudent(ctx echo.Context) error {
	var req views.MeLogin //TODO temp make separate user creation struct
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errz.NewBadRequest("invalid request"))
	}
	if strings.TrimSpace(req.Password) == "" {
		return ctx.JSON(http.StatusBadRequest, errz.NewBadRequest("password is required"))
	}
	if req.Username == nil || strings.TrimSpace(*req.Username) == "" {
		return ctx.JSON(http.StatusBadRequest, errz.NewBadRequest("username is required"))
	}
	if req.Email == nil || strings.TrimSpace(*req.Email) == "" {
		return ctx.JSON(http.StatusBadRequest, errz.NewBadRequest("email is required"))
	}

	cc := ctx.(*CustomContext)
	err = h.s.RegisterStudent(req, cc.user)
	return errz.HandleErrx(ctx, err)
}
