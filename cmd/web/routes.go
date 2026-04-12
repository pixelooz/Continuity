package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (b *backend) setupRoutes(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(b.requestLogger())
	e.Use(b.csrf())

	e.Static("/static", "ui/static")

	auth := e.Group("/user")
	auth.GET("/signup", b.viewUserSignupForm)
	auth.POST("/signup", b.postUserSignupForm)
	auth.GET("/demo", b.demoView)

	auth.GET("/login", b.viewUserLoginForm)
	auth.POST("/login", b.postUserLoginForm)

	auth.POST("/logout", b.postUserLogoutForm)

	protected := e.Group("/v1")
	protected.Use(b.authenticate())

	protected.GET("/home", b.recentNotesView)
	// todo: /home/newest, b.newestNotesView
	// todo: /home/oldest, b.oldestNotesView
	protected.GET("/notes", b.notesView)

	protected.GET("/collection", b.viewRootCollection)
	protected.POST("/collection/:collection-id", b.postCreateCollectionForm)
	protected.GET("/collection/:collection-id", b.viewCollectionPage)
	protected.POST("/collection/:collection-id/delete", b.deleteCollectionPage)

	protected.GET("/note/create", b.viewRootCreateNoteForm)
	protected.GET("/note/create/:collection-id", b.viewCreateNoteForm)
	protected.POST("/note/create/:collection-id", b.postCreateNoteForm)

	protected.GET("/note/:note-id", b.viewNotePage)
	protected.POST("/note/:note-id/delete", b.deleteNotePage)
}
