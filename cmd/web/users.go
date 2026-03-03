package main

import (
	"Continuity/internal/validate"
	"net/http"

	"github.com/labstack/echo/v4"
)

// A UserView is a basic representation of the user.
type UserView struct {
	Name     string
	Username string
	Email    string
}

// userSignupForm is what gets passed to the HTML for data injection.
type userSignupForm struct {
	Name               string `form:"name"`
	Username           string `form:"username"`
	Email              string `form:"email"`
	Password           string `form:"password"`
	validate.Validator `form:"_"`
}

func (b *backend) userSignupFormView(c echo.Context) error {
	td := b.newTemplateData(c)
	td.Form = new(userSignupForm)
	if f, ok := td.Form.(*userSignupForm); ok {
		f.FieldErrors = make(map[string]string)
		f.AddFieldError("name", "some error")
	}
	return c.Render(http.StatusOK, "user_signup.gohtml", td)
}

type userLoginForm struct {
	Username           string `form:"username"`
	Password           string `form:"password"`
	validate.Validator `form:"_"`
}

func (b *backend) userLoginFormView(c echo.Context) error {
	td := b.newTemplateData(c)
	td.Form = new(userLoginForm)
	if f, ok := td.Form.(*userLoginForm); ok {
		f.FieldErrors = make(map[string]string)
		f.AddFieldError("name", "some error")
	}
	return c.Render(http.StatusOK, "user_login.gohtml", td)
}
