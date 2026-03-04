package main

import (
	"Continuity/internal/data"
	"Continuity/internal/validate"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// authenticate extracts the user from the session cookie and sets it in the context.
func (b *backend) authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				b.clearSessionCookie(c)
				return c.Redirect(http.StatusSeeOther, "/user/login")
			}
			token := cookie.Value

			v := validate.NewValidator()
			if validate.ValidatePlainToken(v, token); !v.Valid() {
				return c.String(http.StatusUnauthorized, "invalid token from validation")
			}
			user, err := b.models.Users.GetForToken(
				c.Request().Context(),
				data.ScopeAuthentication, token,
			)
			if err != nil {
				b.clearSessionCookie(c)
				switch {
				case errors.Is(err, data.ErrRecordNotFound):
					return c.Redirect(http.StatusSeeOther, "/user/login")
				default:
					return c.String(http.StatusInternalServerError, "internal server error")
				}
			}
			b.contextSetUser(c, user)
			return next(c)
		}
	}
}

// requestLogger logs every request that goes through this middleware primarily for
// recording latency of several routes.
func (b *backend) requestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()
			err := next(c)
			dur := time.Since(start)

			b.zlog.Info().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Int("status", res.Status).
				Str("ip", c.RealIP()).
				Dur("duration", dur).Msg("")

			return err
		}
	}
}

// csrf configures CSRF middleware to restrict stateful requests from sites other
// than our own. Reduces the threat of CSRF attacks significantly in general.
func (b *backend) csrf() echo.MiddlewareFunc {
	return middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "form:csrf_token",
		CookieName:     "csrf_token",
		CookieHTTPOnly: true,
		CookieSecure:   b.conf.env != "dev",
		CookieSameSite: http.SameSiteLaxMode,
	})
}
