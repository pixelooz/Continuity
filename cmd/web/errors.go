package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Error struct {
	ErrType string
	ErrCode int
	Message string
}

// renderBadRequestErr renders an unauthorized error page.
func (b *backend) renderUnauthorizedErr(c echo.Context, msg string) error {
	return c.Render(http.StatusUnauthorized, "perma_error.gohtml", Error{
		ErrType: "Unauthorized",
		ErrCode: 401, Message: msg,
	})
}

// renderBadRequestErr renders a bad request error page.
func (b *backend) renderBadRequestErr(c echo.Context, msg string) error {
	return c.Render(http.StatusBadRequest, "perma_error.gohtml", Error{
		ErrType: "Bad Request",
		ErrCode: 400, Message: msg,
	})
}

// renderNotFoundErr renders a bad request error page.
func (b *backend) renderNotFoundErr(c echo.Context, msg string) error {
	return c.Render(http.StatusNotFound, "perma_error.gohtml", Error{
		ErrType: "Not Found",
		ErrCode: 404, Message: msg,
	})
}

// renderInternalServerErr renders a internal server error page.
func (b *backend) renderInternalServerErr(c echo.Context) error {
	return c.Render(http.StatusInternalServerError, "perma_error.gohtml", Error{
		ErrType: "Internal Server Error",
		ErrCode: 500,
		Message: "something went wrong",
	})
}
