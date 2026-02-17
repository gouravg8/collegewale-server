package handlers

import (
	"collegeWaleServer/internal/db"
	"collegeWaleServer/internal/enums/roles"
	"collegeWaleServer/internal/model"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CustomContext struct {
	echo.Context
	user *model.User
}

func (c *CustomContext) User() *model.User {
	if c == nil {
		return nil
	}
	return c.user
}

func (c *CustomContext) SetUser(user *model.User) {
	c.user = user
}

func (c *CustomContext) HasRole(role roles.Roles) bool {
	if c == nil || c.user == nil {
		return false
	}
	for _, r := range c.User().Roles {
		if role == r.Name {
			return true
		}
	}
	return false
}

func WithRole(f echo.HandlerFunc, roles ...roles.Roles) echo.HandlerFunc {
	return func(c echo.Context) error {

		cc := c.(*CustomContext)
		allowed := false
		for _, role := range roles {
			if cc.HasRole(role) {
				allowed = true
			}
		}

		if !allowed {
			return c.String(http.StatusForbidden, "You are not allowed to access this resource")
		}
		return f(c)
	}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &CustomContext{Context: c}
		token := c.Get("user").(*jwt.Token)
		if token == nil {
			return next(cc)
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			return next(cc)
		}
		//TODO make cache package to always have a user cache and minimize db calls
		//cachedUser, found := cache.UsersCache.Get(sub)
		//if found {
		//	cc.user = &cachedUser
		//	return next(cc)
		//}
		var dbuser model.User
		err = db.DB.Model(&model.User{}).Preload("Roles").Where("username = ?", sub).First(&dbuser).Error
		if err != nil {
			log.Error("error fetching user details", err)
			return next(cc)
		}
		//cache.UsersCache.Put(sub, dbuser)
		cc.user = &dbuser
		return next(cc)
	}
}
