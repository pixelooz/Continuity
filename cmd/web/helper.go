package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/labstack/echo/v4"
)

// decodePostForm decodes the POST form into the given dst struct.
func (b *backend) decodePostForm(c echo.Context, dst any) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	err := b.decoder.Decode(dst, c.Request().PostForm)
	if err != nil {
		decErr, ok := errors.AsType[*form.InvalidDecoderError](err)
		if ok {
			panic(decErr.Error())
		}
		return fmt.Errorf("couldn't decode form: %w", err)
	}
	return nil
}

// renderWithFieldErr is a convenience method that re-renders the named template
// with the same form data. It's intended to be used after doing some sort of
// validation where the form now contains some validation field errors that will
// be rendered with http.StatusUnprocessableEntity status code.
// Params: (name) of the template, (form) is the pointer object that now contains
// the field errors and the previous form's data.
func (b *backend) renderWithFieldErr(c echo.Context, name string, form any) error {
	pd := b.NewPageData(c)
	pd.Form = form
	return c.Render(http.StatusUnprocessableEntity, name, pd)
}

// renderErrFormData works the same way renderWithFieldErr does but for pages
// that are not pure forms but also have some data on them.
func (b *backend) renderErrFormData(c echo.Context, name string, form, data any) error {
	pd := b.NewPageData(c)
	pd.Form = form
	pd.Data = data
	return c.Render(http.StatusUnprocessableEntity, name, pd)
}

func (b *backend) clearSessionCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})
}

// isAuthenticated returns true if the user is authenticated.
func (b *backend) isAuthenticated(c echo.Context) bool {
	return c.Get("user") != nil
}
