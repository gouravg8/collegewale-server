package handlers

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums"
	service "collegeWaleServer/internal/service"
	"collegeWaleServer/internal/utils/common"
	"collegeWaleServer/internal/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StudentHandler struct {
	st *service.StudentService
}

func NewStudentHandler(group *echo.Group, st *service.StudentService) *StudentHandler {
	h := &StudentHandler{st: st}
	group.POST("/student/list", h.ListStudents)
	group.PUT("/student/update-status", h.UpdateStudentStatus)
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
	res, err := h.st.ListStudents(filter)
	if err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

func (h *StudentHandler) UpdateStudentStatus(c echo.Context) error {
	id := c.QueryParam("id")
	st := c.QueryParam("status")
	maskedId := common.MaskedId(id)
	status := enums.StudentStatus(st)
	if err := status.IsValid(); err != nil {
		return err
	}
	if err := h.st.UpdateStudentStatus(maskedId, status); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, "status updated successfully!")
}
