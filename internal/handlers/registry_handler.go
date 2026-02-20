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
	group.POST("/register/college/user", WithRole(h.RegisterCollegeAccount, roles.Admin))
	group.POST("/register/student", WithRole(h.RegisterStudent, roles.Admin, roles.College))
	return h
}

func (h Registry) RegisterCollege(c echo.Context) error {
	var req views.College
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("invalid request"))
	}
	if err := req.IsValidRequest(); err != nil {
		return errz.HandleErrx(c, err)
	}
	cc := c.(*CustomContext)
	if err := h.s.RegisterCollege(req, cc.user); err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, views.Response{Message: "success"})
}

func (h Registry) RegisterStudent(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc == nil {
		return c.JSON(http.StatusOK, errz.NewBadRequest("user not found."))
	}
	var req views.StudentForm
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("invalid request"))
	}
	err = req.IsValid()
	if err != nil {
		return errz.HandleErrx(c, err)
	}
	err = h.s.RegisterStudent(req, cc.user)
	return errz.HandleErrx(c, err)
}

func (h Registry) RegisterCollegeAccount(ctx echo.Context) error {
	cc := ctx.(*CustomContext)
	if cc == nil {
		return ctx.JSON(http.StatusOK, errz.NewBadRequest("user not found."))
	}
	var req views.MeLogin //TODO TEMP make separate College user creation struct view
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

	err = h.s.RegisterCollegeAccount(req, cc.user)
	return errz.HandleErrx(ctx, err)
}
