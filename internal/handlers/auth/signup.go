package auth_handler

import (
	"collegeWaleServer/internal/models"
	service "collegeWaleServer/internal/services/auth"
	"collegeWaleServer/internal/views"
	auth_view "collegeWaleServer/internal/views/auth"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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
	group.POST("/set-password", h.SetPassword)
	group.POST("/college-login", h.CollegeLogin)
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

	if req.Email == "" && req.Code == "" {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: "Enter Email or Code"})
	}

	if req.Password == "" || req.ConfirmPassword == "" {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: "Password is required"})
	}

	if req.Password != req.ConfirmPassword {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: "Both Password must match"})
	}

	if err := h.authService.SetPassword(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, views.Response{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, views.Response{Message: "success"})
}

func (h AuthHandler) CollegeLogin(c echo.Context) error {
	var req auth_view.CollegeLogin
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, views.Response{Message: "Can not map the request"})
	}

	if req.Code == "" && req.Email == "" {
		return c.JSON(http.StatusBadRequest, views.Response{Message: "College Code or Email is required"})
	}

	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, views.Response{Message: "Password is required"})
	}

	college, err := h.authService.CollegeLogin(req)

	token, err := h.generateToken(college)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Response{Message: err.Error()})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Response{
			Data: auth_view.CollegeLoginResponse{
				Name:  college.Name,
				Code:  college.Code,
				Email: college.Email,
				Token: token,
			},
		})
	}
}

func (h AuthHandler) generateToken(college *models.College) (string, error) {
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
		return "", err
	}
	return t, nil
}
