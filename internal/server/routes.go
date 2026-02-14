package server

import (
	"collegeWaleServer/internal/database"
	auth_handler "collegeWaleServer/internal/handlers/auth"
	service "collegeWaleServer/internal/services/auth"
	"collegeWaleServer/internal/views"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	dbService := database.New()

	authGroup := e.Group("/auth")
	apiV1Group := e.Group("/api/v1")
	apiV1Group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
	}))

	authService := service.NewAuthService(dbService.DB)

	auth_handler.NewAuthHandler(authGroup, authService)

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{allowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	e.GET("/", s.HelloWorldHandler)
	e.GET("/health", s.healthHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, views.Response{
		Status:  http.StatusOK,
		Message: "Welcome to CollegeWale API",
		Data: map[string]string{
			"version": "1.0.0",
			"service": "collegewale-server",
		},
	})
}

func (s *Server) healthHandler(c echo.Context) error {
	healthStatus := s.db.Health()

	status := http.StatusOK
	message := "Service is healthy"

	if healthStatus["status"] == "down" {
		status = http.StatusServiceUnavailable
		message = "Service is unhealthy"
	}

	return c.JSON(status, views.Response{
		Status:  status,
		Message: message,
		Data:    healthStatus,
	})
}
