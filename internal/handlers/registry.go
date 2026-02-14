package handlers

import (
	"collegeWaleServer/errz"
	service "collegeWaleServer/internal/services/auth"
	"collegeWaleServer/internal/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Registry struct {
	s *service.RegistryService
}

func NewRegistryHandler(group *echo.Group, registryService *service.RegistryService) *Registry {
	h := &Registry{
		s: registryService,
	}
	//group.POST("/register/college", WithRole(h.RegisterCollege, roles.Admin))
	group.POST("/register/college", h.RegisterCollege)
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
	if err := h.s.RegisterCollege(req); err != nil {
		return errz.HandleErrx(ctx, err)
	}
	return ctx.JSON(http.StatusOK, views.Response{Message: "success"})
}
