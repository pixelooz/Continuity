package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (b *backend) setupRoutes(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(b.RequestLogger())
	e.Use(CSRF())

	e.Static("/static", "ui/static")

	pub := e.Group("/v1")
	pub.GET("/home", b.homeView)
	pub.GET("/notes", b.notesView)

	auth := e.Group("/user")
	auth.GET("/signup", b.userSignupFormView)
	auth.POST("/signup", b.userSignupFormPost)

	auth.GET("/login", b.userLoginFormView)
	auth.POST("/login", b.userLoginFormPost)
}
