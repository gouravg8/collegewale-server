package auth_handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"collegeWaleServer/internal/models"
	service "collegeWaleServer/internal/services/auth"
	"collegeWaleServer/internal/views"
	auth_view "collegeWaleServer/internal/views/auth"
)

type AuthHandler struct {
	jwtKey      string
	authService *service.AuthService
}

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func NewAuthHandler(group *echo.Group, authService *service.AuthService) *AuthHandler {
	h := &AuthHandler{
		authService: authService,
		jwtKey:      os.Getenv("JWT_SECRET_KEY"),
	}

	group.POST("/college-signup", h.DoSignup)
	group.POST("/verification", h.Verification)
	group.POST("/college-set-password", h.SetPassword)
	group.POST("/college-login", h.CollegeLogin)
	return h
}

func (h AuthHandler) DoSignup(ctx echo.Context) error {
	var req auth_view.CollegeSignup
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
	}

	college, msg, err := h.authService.CollegeSignup(req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Signup failed",
			Error:   err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, views.Response{
		Status:  http.StatusCreated,
		Message: msg,
		Data: map[string]any{
			"code":  college.Code,
			"email": college.Email,
			"name":  college.Name,
		},
	})
}

func (h AuthHandler) Verification(ctx echo.Context) error {
	token := ctx.QueryParam("token")

	if token == "" {
		return ctx.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Token is required",
			Error:   "missing token parameter",
		})
	}

	college, err := h.authService.GetCollegeByToken(token)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, views.Response{
			Status:  http.StatusUnauthorized,
			Message: "Token verification failed",
			Error:   err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, views.Response{
		Status:  http.StatusOK,
		Message: "Token verified successfully",
		Data: map[string]any{
			"code":  college.Code,
			"email": college.Email,
		},
	})
}

func (h AuthHandler) SetPassword(ctx echo.Context) error {
	var req auth_view.SetPassword
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
	}

	// Validate input
	if req.Email == "" && req.Code == "" {
		return ctx.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Error:   "either email or code must be provided",
		})
	}

	if req.Password == "" || req.ConfirmPassword == "" {
		return ctx.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Error:   "password and confirm password are required",
		})
	}

	if req.Password != req.ConfirmPassword {
		return ctx.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Error:   "passwords do not match",
		})
	}

	if len(req.Password) < 8 {
		return ctx.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Error:   "password must be at least 8 characters long",
		})
	}

	err = h.authService.SetPassword(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, views.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to set password",
			Error:   err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, views.Response{
		Status:  http.StatusOK,
		Message: "Password set successfully",
	})
}

func (h AuthHandler) CollegeLogin(c echo.Context) error {
	var req auth_view.CollegeLogin
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
	}

	// Validate input
	if req.Code == "" && req.Email == "" {
		return c.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Error:   "either college code or email is required",
		})
	}

	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, views.Response{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Error:   "password is required",
		})
	}

	college, err := h.authService.CollegeLogin(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, views.Response{
			Status:  http.StatusUnauthorized,
			Message: "Authentication failed",
			Error:   err.Error(),
		})
	}

	token, err := h.generateToken(college)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Response{
			Status:  http.StatusInternalServerError,
			Message: "Token generation failed",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, views.Response{
		Status:  http.StatusOK,
		Message: "Login successful",
		Data: auth_view.CollegeLoginResponse{
			Name:  college.Name,
			Code:  college.Code,
			Email: college.Email,
			Token: token,
		},
	})
}

func (h AuthHandler) generateToken(college *models.College) (string, error) {
	if h.jwtKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY not configured")
	}

	claims := &jwtCustomClaims{
		Name: college.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "collegewale",
			Subject:   college.Email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(h.jwtKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return t, nil
}
