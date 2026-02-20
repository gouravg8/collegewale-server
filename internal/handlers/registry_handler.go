package handlers

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums/roles"
	service "collegeWaleServer/internal/service/auth"
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
	group.POST("/register/college", WithRole(h.RegisterCollege, roles.Admin))
	group.POST("/register/college/user", WithRole(h.RegisterCollegeAccount, roles.Admin))
	group.POST("/register/student", WithRole(h.RegisterStudent, roles.Admin, roles.College))
	return h
}

func (h Registry) RegisterCollege(c echo.Context) error {
	var req views.CollegeRequest
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
	var req views.CollegeSignup
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errz.NewBadRequest("invalid request"))
	}

	if err = req.IsValid(); err != nil {
		return errz.HandleErrx(ctx, err)
	}

	err = h.s.RegisterCollegeAccount(req, cc.user)
	return errz.HandleErrx(ctx, err)
}
