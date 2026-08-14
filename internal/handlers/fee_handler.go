package handlers

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums/roles"
	"collegeWaleServer/internal/service"
	"collegeWaleServer/internal/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FeeHandler struct {
	fs *service.FeeService
}

func NewFeeHandler(g *echo.Group, fs *service.FeeService) *FeeHandler {
	h := &FeeHandler{fs: fs}
	g.POST("/fees", WithRole(h.CreateFeeRecord, roles.CollegeAdmin))
	g.GET("/fees", WithRole(h.ListFeeRecords, roles.CollegeAdmin))
	g.POST("/fees/payment", WithRole(h.RecordPayment, roles.CollegeAdmin))
	g.GET("/fees/mine", WithRole(h.ListMyFeeRecords, roles.Student))
	return h
}

func (h FeeHandler) CreateFeeRecord(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc == nil || cc.User() == nil {
		return errz.HandleErrx(c, errz.NewUnauthorized("user not found"))
	}
	user := cc.User()
	if user.CollegeID == nil {
		return errz.HandleErrx(c, errz.NewBadRequest("user is not linked to a college"))
	}

	var req views.CreateFeeRecordRequest
	if err := c.Bind(&req); err != nil {
		return errz.HandleErrx(c, errz.NewBadRequest("invalid request"))
	}
	if err := req.IsValid(); err != nil {
		return errz.HandleErrx(c, err)
	}

	res, err := h.fs.CreateFeeRecord(req, *user.CollegeID, user.ID)
	if err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

func (h FeeHandler) ListFeeRecords(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc == nil || cc.User() == nil {
		return errz.HandleErrx(c, errz.NewUnauthorized("user not found"))
	}
	user := cc.User()
	if user.CollegeID == nil {
		return errz.HandleErrx(c, errz.NewBadRequest("user is not linked to a college"))
	}

	res, err := h.fs.ListFeeRecords(*user.CollegeID)
	if err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

func (h FeeHandler) RecordPayment(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc == nil || cc.User() == nil {
		return errz.HandleErrx(c, errz.NewUnauthorized("user not found"))
	}
	user := cc.User()
	if user.CollegeID == nil {
		return errz.HandleErrx(c, errz.NewBadRequest("user is not linked to a college"))
	}

	var req views.RecordPaymentRequest
	if err := c.Bind(&req); err != nil {
		return errz.HandleErrx(c, errz.NewBadRequest("invalid request"))
	}
	if err := req.IsValid(); err != nil {
		return errz.HandleErrx(c, err)
	}

	res, err := h.fs.RecordPayment(req, *user.CollegeID, user.ID)
	if err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

func (h FeeHandler) ListMyFeeRecords(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc == nil || cc.User() == nil {
		return errz.HandleErrx(c, errz.NewUnauthorized("user not found"))
	}
	user := cc.User()
	if user.Student == nil {
		return errz.HandleErrx(c, errz.NewBadRequest("user is not linked to a student profile"))
	}

	res, err := h.fs.ListFeeRecordsForStudent(user.Student.ID)
	if err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, res)
}
