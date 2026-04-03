package handlers

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums/roles"
	service "collegeWaleServer/internal/service/auth"
	"collegeWaleServer/internal/storage"
	"collegeWaleServer/internal/views"
	"collegeWaleServer/internal/views/common"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
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
	file, err := c.FormFile("logo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("college logo is required."))
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	metadata := c.FormValue("metadata")
	var req views.CollegeRequest
	if err := json.Unmarshal([]byte(metadata), &req); err != nil {
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("invalid request"))
	}
	if err := req.IsValidRequest(); err != nil {
		return errz.HandleErrx(c, err)
	}
	cc := c.(*CustomContext)
	if cc == nil {
		return c.JSON(http.StatusOK, errz.NewBadRequest("user not found."))
	}
	//objectKey := fmt.Sprintf("logos/%s/%s-%s", req.Code, "a", "test.png")
	objectKey := "test.png"
	client := storage.InitR2Client()
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String("collegewala-server"),
		Key:         aws.String(objectKey),
		Body:        src,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		log.Errorf("Failed to upload file: %v", err)
		return c.JSON(http.StatusBadRequest, errz.NewBadRequest("upload file failed."))
	}

	if err := h.s.RegisterCollege(req, cc.user); err != nil {
		return errz.HandleErrx(c, err)
	}
	return c.JSON(http.StatusOK, view.Response{Message: "success"})
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
