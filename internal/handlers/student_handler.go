package handlers

import (
	"collegeWaleServer/errz"
	service "collegeWaleServer/internal/service"
	"collegeWaleServer/internal/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StudentHandler struct {
	st *service.StudentService
}

func NewStudentHandler(g *echo.Group, st *service.StudentService) *StudentHandler {
	h := &StudentHandler{st: st}
	g.POST("/students", h.ListStudents)
	return h
}

// ListStudents retrieves all students with optional filtering
func (h *StudentHandler) ListStudents(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc == nil {
		return errz.HandleErrx(c, errz.NewUnauthorized("user not found"))
	}
	var filter views.StudentFilter
	err := c.Bind(&filter)
	if err != nil {
		return errz.NewBadRequest("invalid request :: failed to bind request")
	}
	res, err := h.st.ListStudents(cc.user, filter)
	if err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, res)
}
