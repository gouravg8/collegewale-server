package handlers

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums/roles"
	appService "collegeWaleServer/internal/service/application"
	"collegeWaleServer/internal/views"
	appViews "collegeWaleServer/internal/views/application"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ApplicationHandler struct {
	s *appService.ApplicationService
}

func NewApplicationHandler(group *echo.Group, svc *appService.ApplicationService) *ApplicationHandler {
	h := &ApplicationHandler{s: svc}
	group.POST("/applications", WithRole(h.Create, roles.Admin, roles.College))
	group.GET("/applications", WithRole(h.List, roles.Admin, roles.College))
	group.GET("/applications/:id", WithRole(h.Get, roles.Admin, roles.College, roles.Student))
	group.PUT("/applications/:id/status", WithRole(h.UpdateStatus, roles.Admin, roles.College))
	group.GET("/my/applications", WithRole(h.MyApplications, roles.Student))
	return h
}

func (h *ApplicationHandler) Create(c echo.Context) error {
	cc := c.(*CustomContext)
	var req appViews.CreateApplicationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("invalid request"))
	}
	if err := req.IsValid(); err != nil {
		return errz.HandleErrx(c, err)
	}

	app, err := h.s.CreateApplication(req, cc.user)
	if err != nil {
		return errz.HandleErrx(c, err)
	}

	resp := appViews.NewApplicationResponse(app)
	return c.JSON(http.StatusCreated, views.Response{
		Status:  http.StatusCreated,
		Message: "application created",
		Data:    resp,
	})
}

func (h *ApplicationHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("invalid application id"))
	}

	app, err := h.s.GetApplication(uint(id))
	if err != nil {
		return errz.HandleErrx(c, err)
	}

	resp := appViews.NewApplicationResponse(app)
	return c.JSON(http.StatusOK, views.Response{Data: resp})
}

func (h *ApplicationHandler) List(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc.user.CollegeID == nil {
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("user is not associated with a college"))
	}

	status := c.QueryParam("status")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	apps, total, err := h.s.ListApplications(*cc.user.CollegeID, status, page, pageSize)
	if err != nil {
		return errz.HandleErrx(c, err)
	}

	var respList []appViews.ApplicationResponse
	for i := range apps {
		respList = append(respList, appViews.NewApplicationResponse(&apps[i]))
	}

	return c.JSON(http.StatusOK, views.ListResponse{
		TotalRecords: int(total),
		Data:         respList,
	})
}

func (h *ApplicationHandler) UpdateStatus(c echo.Context) error {
	cc := c.(*CustomContext)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("invalid application id"))
	}

	var req appViews.UpdateStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("invalid request"))
	}
	if err := req.IsValid(); err != nil {
		return errz.HandleErrx(c, err)
	}

	app, err := h.s.UpdateStatus(uint(id), req, cc.user)
	if err != nil {
		return errz.HandleErrx(c, err)
	}

	resp := appViews.NewApplicationResponse(app)
	return c.JSON(http.StatusOK, views.Response{
		Message: "status updated",
		Data:    resp,
	})
}

func (h *ApplicationHandler) MyApplications(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc.user.Student == nil {
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("no student profile linked"))
	}

	apps, err := h.s.GetStudentApplications(cc.user.Student.ID)
	if err != nil {
		return errz.HandleErrx(c, err)
	}

	var respList []appViews.ApplicationResponse
	for i := range apps {
		respList = append(respList, appViews.NewApplicationResponse(&apps[i]))
	}

	return c.JSON(http.StatusOK, views.ListResponse{
		TotalRecords: len(respList),
		Data:         respList,
	})
}
