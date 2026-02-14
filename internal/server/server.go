package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"collegeWaleServer/internal/database"
)

type Server struct {
	e  *echo.Echo
	db database.Service
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	e := echo.New()
	s.e = e
	s.db = *database.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	/*----echo-config----*/
	e.Server.Handler = s.RegisterRoutes()
	e.Server.IdleTimeout = time.Minute

	return nil
}

func (s *Server) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	s.e.Logger.Fatal(s.e.Start(fmt.Sprintf(":%s", port)))
}

func (s *Server) dbHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
