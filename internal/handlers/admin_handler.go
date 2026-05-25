package handlers

import (
	"collegeWaleServer/internal/service"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	as *service.AdminService
}

func NewAdminHandler(g *echo.Group, as *service.AdminService) *AdminHandler {
	h := &AdminHandler{as: as}
	return h

}
