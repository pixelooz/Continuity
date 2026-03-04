package main

import (
	"Continuity/internal/data"

	"github.com/labstack/echo/v4"
)

// contextSetUser sets the user in echo's context, it's merely a wrapper around c.Set().
func (b *backend) contextSetUser(c echo.Context, user *data.User) {
	c.Set("user", user)
}

// contextGetUser returns the user from echo's context or panics if it's not found.
func (b *backend) contextGetUser(c echo.Context) *data.User {
	user, ok := c.Get("user").(*data.User)
	if !ok {
		panic("user not found in context")
	}
	return user
}
