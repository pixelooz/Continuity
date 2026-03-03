package main

import (
	"Continuity/internal/data"
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

// userSignupForm holds the data that we get from the signup form.
type userSignupForm struct {
	Name               string `form:"name"`
	Username           string `form:"username"`
	Email              string `form:"email"`
	Password           string `form:"password"`
	validate.Validator `form:"_"`
}

// userSignupFormView renders the user signup form.
func (b *backend) userSignupFormView(c echo.Context) error {
	pd := b.NewPageData(c)
	pd.Form = new(userSignupForm)
	return c.Render(http.StatusOK, "user_signup.gohtml", pd)
}

// userSignupFormPost handles signup requests and validates the data before registering
// the user, sending the appropriate error, if any.
func (b *backend) userSignupFormPost(c echo.Context) error {
	var signupForm userSignupForm
	if err := b.decodePostForm(c, &signupForm); err != nil {
		return c.String(http.StatusBadRequest, "invalid form data")
	}
	user := &data.User{
		Name:      signupForm.Name,
		Username:  signupForm.Username,
		Email:     signupForm.Email,
		Activated: false,
	}
	err := user.Password.SetPass(signupForm.Password)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	validate.ValidateUser(&signupForm.Validator, user)
	if !signupForm.Validator.Valid() {
		pd := b.NewPageData(c)
		pd.Form = &signupForm
		return c.Render(http.StatusUnprocessableEntity, "user_signup.gohtml", pd)
	}
	return c.Redirect(http.StatusSeeOther, "/user/login")
}

// userLoginForm holds the data that we get from the login form.
type userLoginForm struct {
	Username           string `form:"username"`
	Password           string `form:"password"`
	validate.Validator `form:"_"`
}

// userLoginFormView renders the user login form.
func (b *backend) userLoginFormView(c echo.Context) error {
	pd := b.NewPageData(c)
	pd.Form = new(userLoginForm)
	return c.Render(http.StatusOK, "user_login.gohtml", pd)
}

// userLoginFormPost handles login requests and validates the data before allowing
// the user to log in, sending the appropriate error, if any.
func (b *backend) userLoginFormPost(c echo.Context) error {
	var loginForm userLoginForm
	if err := b.decodePostForm(c, &loginForm); err != nil {
		return c.String(http.StatusBadRequest, "invalid form data")
	}
	validate.ValidateUsername(&loginForm.Validator, loginForm.Username)
	validate.ValidatePlainPassword(&loginForm.Validator, loginForm.Password)
	if !loginForm.Validator.Valid() {
		pd := b.NewPageData(c)
		pd.Form = &loginForm
		return c.Render(http.StatusUnprocessableEntity, "user_login.gohtml", pd)
	}
	return c.Redirect(http.StatusSeeOther, "/v1/notes")
}
