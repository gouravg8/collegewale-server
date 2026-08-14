package handlers

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums/roles"
	"collegeWaleServer/internal/service"
	"collegeWaleServer/internal/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CollegeHandler struct {
	cs  *service.CollegeService
	css *service.CourseService
}

func NewCollegeHandler(g *echo.Group, cs *service.CollegeService, css *service.CourseService) *CollegeHandler {
	h := &CollegeHandler{cs: cs, css: css}
	g.POST("/college/courses", h.ListCourses)
	g.GET("/college/stats", h.GetStats)
	g.GET("/colleges", WithRole(h.ListColleges, roles.Admin))
	return h
}

func (h CollegeHandler) ListColleges(c echo.Context) error {
	res, err := h.cs.ListColleges()
	if err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

func (h CollegeHandler) ListCourses(c echo.Context) error {
	var filter views.CoursesFilter
	if err := c.Bind(&filter); err != nil {
		return errz.HandleErrx(c, errz.NewBadRequest("invalid request :: failed to bind request"))
	}
	res, err := h.css.ListCourses(filter)
	if err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

func (h CollegeHandler) GetStats(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc == nil {
		return errz.HandleErrx(c, errz.NewUnauthorized("user not found"))
	}
	res, err := h.cs.GetStats(cc.user)
	if err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, res)
}
