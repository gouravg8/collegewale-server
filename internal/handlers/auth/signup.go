package auth_handler

import (
	service "collegeWaleServer/internal/services/auth"
	"collegeWaleServer/internal/views"
	auth_view "collegeWaleServer/internal/views/auth"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	jwtKey      string
	authService *service.AuthService
}

func NewAuthHandler(group *echo.Group, authService *service.AuthService) *AuthHandler {
	h := &AuthHandler{
		authService: authService,
		jwtKey:      os.Getenv("JWT_SECRET_KEY"),
	}

	group.POST("/college-signup", h.DoSignup)
	return h
}

func (h AuthHandler) DoSignup(ctx echo.Context) error {
	var req auth_view.CollegeSignup
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: "can not map", Data: err})
	}

	_, err = h.authService.CollegeSignup(req)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, views.Response{Message: err.Error()})
	}
	fmt.Println("hi from do signup handler")
	return ctx.JSON(http.StatusOK, views.Response{Message: "success"})
}
