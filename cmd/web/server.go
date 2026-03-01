package main

import (
	"Continuity/internal/data"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// config struct holds the options necessary for a successful initialization
// of the backend struct.
type config struct {
	env  string
	addr string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	smtp struct {
		sender   string
		host     string
		port     int
		username string
		password string
	}
	cors struct {
		allowedOrigins []string
	}
}

// backend struct holds the necessary dependencies for the backend to run.
type backend struct {
	renderer *TemplateRenderer
	conf     config
	model    *data.Models
	wg       sync.WaitGroup
	zlog     *zerolog.Logger
}

// serve sets up the server configuration and starts the server on the configured
// address and spins up a background goroutine for graceful shutdown in case of a
// SigInt/SigTerm signal.
func (b *backend) serve() error {
	e := echo.New()
	defer func() { _ = e.Close() }()

	renderer, err := NewCacheRenderer()
	if err != nil {
		return err
	}
	b.renderer = renderer
	e.Renderer = renderer
	b.setupRoutes(e)
	srv := http.Server{
		Handler: e, Addr: b.conf.addr,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       90 * time.Second,
	}
	shutdownErr := make(chan error, 1)
	go func() {
		q := make(chan os.Signal, 1)
		signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
		s := <-q
		b.zlog.Info().
			Str("signal", s.String()).
			Msg("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = srv.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}
		b.zlog.Info().
			Str("address", srv.Addr).
			Msg("completing background tasks")

		b.wg.Wait()
		shutdownErr <- nil
	}()
	b.zlog.Info().
		Str("env", b.conf.env).
		Str("addr", srv.Addr).Msg("starting server")

	if err = srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("couldn't listen on addr=%s: %w", srv.Addr, err)
		}
	}
	if err = <-shutdownErr; err != nil {
		return err
	}
	b.zlog.Info().
		Str("env", b.conf.env).
		Str("addr", srv.Addr).Msg("server stopped")
	return nil
}
