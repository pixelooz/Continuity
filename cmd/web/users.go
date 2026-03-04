package main

import (
	"Continuity/internal/data"
	"Continuity/internal/validate"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	ValidOneDay   = 24 * time.Hour
	ValidOneWeek  = 7 * 24 * time.Hour
	ValidOneMonth = 30 * 24 * time.Hour
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

// viewUserSignupForm renders the user signup form.
func (b *backend) viewUserSignupForm(c echo.Context) error {
	pd := b.NewPageData(c)
	pd.Form = new(userSignupForm)
	return c.Render(http.StatusOK, "user_signup.gohtml", pd)
}

// postUserSignupForm handles signup requests and validates the data before registering
// the user, sending the appropriate error, if any.
func (b *backend) postUserSignupForm(c echo.Context) error {
	var signupForm userSignupForm

	if err := b.decodePostForm(c, &signupForm); err != nil {
		return b.renderBadRequestErr(c, "invalid form data")
	}
	user := &data.User{ID: uuid.New(),
		Name:      signupForm.Name,
		Username:  signupForm.Username,
		Activated: false,
		Email:     signupForm.Email,
	}
	err := user.Password.SetPass(signupForm.Password)
	if err != nil {
		return b.renderInternalServerErr(c)
	}
	validate.ValidateUser(&signupForm.Validator, user)
	if !signupForm.Validator.Valid() {
		return b.renderWithFieldErr(c, "user_signup.gohtml", &signupForm)
	}
	err = b.models.Users.Insert(c.Request().Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			signupForm.AddFieldError("email", "email already in use")
			return b.renderWithFieldErr(c, "user_signup.gohtml", &signupForm)
		case errors.Is(err, data.ErrDuplicateUsername):
			signupForm.AddFieldError("username", "username already in use")
			return b.renderWithFieldErr(c, "user_signup.gohtml", &signupForm)
		default:
			b.zlog.Err(err).
				Str("handler", "postUserSignupForm").
				Msg("failed to insert user into db")
			return b.renderInternalServerErr(c)
		}
	}
	token, err := b.models.Tokens.NewToken(c.Request().Context(), user.ID, ValidOneDay, data.ScopeAuthentication)
	if err != nil {
		b.zlog.Err(err).
			Str("handler", "postUserSignupForm").
			Msg("failed to create a new token")
		return b.renderInternalServerErr(c)
	}
	c.SetCookie(&http.Cookie{Name: "session",
		Value:    token.PlainText,
		HttpOnly: true,
		Secure:   b.conf.env != "dev",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(ValidOneDay),
	})
	return c.Redirect(http.StatusSeeOther, "/v1/home")
}

// userLoginForm holds the data that we get from the login form.
type userLoginForm struct {
	Username           string `form:"username"`
	Password           string `form:"password"`
	validate.Validator `form:"_"`
}

// viewUserLoginForm renders the user login form.
func (b *backend) viewUserLoginForm(c echo.Context) error {
	pd := b.NewPageData(c)
	pd.Form = new(userLoginForm)
	return c.Render(http.StatusOK, "user_login.gohtml", pd)
}

// postUserLoginForm handles login requests and validates the data before allowing
// the user to log in, sending the appropriate error, if any.
func (b *backend) postUserLoginForm(c echo.Context) error {
	var loginForm userLoginForm

	if err := b.decodePostForm(c, &loginForm); err != nil {
		return b.renderBadRequestErr(c, "invalid form data")
	}
	validate.ValidateUsername(&loginForm.Validator, loginForm.Username)

	validate.ValidatePassword(&loginForm.Validator, loginForm.Password)
	if !loginForm.Validator.Valid() {
		return b.renderWithFieldErr(c, "user_login.gohtml", &loginForm)
	}
	user, err := b.models.Users.GetByUsername(c.Request().Context(), loginForm.Username)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			loginForm.AddFieldError("username", "username not found")
			return b.renderWithFieldErr(c, "user_login.gohtml", &loginForm)
		}
		return b.renderInternalServerErr(c)
	}
	matches, err := user.Password.CheckPass(loginForm.Password)
	if err != nil {
		b.zlog.Err(err).
			Str("handler", "postUserLoginForm").
			Msg("failed to check password")
		return b.renderInternalServerErr(c)
	}
	if !matches {
		loginForm.AddFieldError("password", "invalid password")
		return b.renderWithFieldErr(c, "user_login.gohtml", &loginForm)
	}
	token, err := b.models.Tokens.NewToken(c.Request().Context(), user.ID, ValidOneMonth, data.ScopeAuthentication)
	if err != nil {
		b.zlog.Err(err).
			Str("handler", "postUserLoginForm").
			Msg("failed to create a new token")
		return b.renderInternalServerErr(c)
	}
	c.SetCookie(&http.Cookie{Name: "session",
		Value:    token.PlainText,
		HttpOnly: true,
		Secure:   b.conf.env != "dev",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(ValidOneMonth),
	})
	return c.Redirect(http.StatusSeeOther, "/v1/notes")
}
