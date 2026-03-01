package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (b *backend) setupRoutes(e *echo.Echo) {
	e.Use(middleware.Gzip())
	e.Static("/", "ui/static")

	prod := e.Group("/v1")
	prod.GET("/home", b.homeView)
	prod.GET("/notes", b.notesView)
}
