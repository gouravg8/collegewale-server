package auth_handler

import (
	service "collegeWaleServer/internal/services/auth"
	"collegeWaleServer/internal/views"
	auth_view "collegeWaleServer/internal/views/auth"
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
	group.POST("/verification", h.Verification)
	group.POST("/set-password", h.SetPassword)
	return h
}

func (h AuthHandler) DoSignup(ctx echo.Context) error {
	var req auth_view.CollegeSignup
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: "can not map", Data: err})
	}

	var msg string
	_, msg, err = h.authService.CollegeSignup(req)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, views.Response{Message: err.Error()})
	}

	if msg != "" {
		return ctx.JSON(http.StatusOK, views.Response{Message: msg})
	} else {
		return ctx.JSON(http.StatusOK, views.Response{Message: "sucess"})
	}
}

func (h AuthHandler) Verification(ctx echo.Context) error {
	token := ctx.QueryParam("token")

	if token == "" {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: "token is required"})
	}

	college, err := h.authService.GetCollegeByToken(token)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, views.Response{Message: "Invalid or expired token"})
	}

	return ctx.JSON(http.StatusOK, views.Response{
		Message: "Token verified, proceed to set password",
		Data:    map[string]any{"college_id": college.ID},
	})
}

func (h AuthHandler) SetPassword(ctx echo.Context) error {
	var req auth_view.SetPassword
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			views.Response{Message: "can not map", Data: err})
	}

	if req.Password != req.ConfirmPassword {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: "Both Password must match"})
	}

	if err := h.authService.SetPassword(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, views.Response{Message: "success"})

}

func (h AuthHandler) Login(ctx echo.Context) error {
	var req auth_view.CollegeLogin
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: "can not map request", Data: err})
	}
	
	if err := h.authService.CollegeLogin(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, views.Response{Message: "success"})
}
