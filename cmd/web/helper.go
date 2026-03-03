package main

import (
	"time"

	"github.com/labstack/echo/v4"
)

func (b *backend) newTemplateData(c echo.Context) *templateData {
	return &templateData{
		CreatedAt: time.Now().Format("03 Jan 2006"),
		CSRFToken: c.Get("csrf").(string),
		// UserView: b.currentUserView(c),
		// IsAuthenticated: b.isAuthenticated(c),
	}
}
