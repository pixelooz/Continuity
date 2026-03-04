package main

import (
	"Continuity/internal/data"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".envrc"); err != nil {
		log.Fatalf("error loading .envrc: %v", err)
	}
	var conf config
	runFlags(&conf)

	zlog := zerolog.New(zerolog.NewConsoleWriter()).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().Caller().Logger()

	db, err := openDb(conf)
	if err != nil {
		zlog.Fatal().Err(err).Msg("failed to open db")
	}
	defer func() {
		if err = db.Close(); err != nil {
			zlog.Err(err).Msg("couldn't close db")
		}
	}()
	zlog.Info().Msg("database connection established")

	err = CreateHighlightCSS("ui/static/css/highlight.css")
	if err != nil {
		zlog.Fatal().Err(err).
			Msg("Failed to write highlight.css")
	}
	bknd := backend{
		decoder: form.NewDecoder(),
		zlog:    &zlog,
		conf:    conf,
		models:  data.NewModels(db),
	}
	if err = bknd.serve(); err != nil {
		zlog.Err(err).
			Str("addr", bknd.conf.addr).
			Msg("couldn't start server")
	}
}

func runFlags(conf *config) {
	flag.StringVar(&conf.env, "env", "dev", "Environment: dev, prod")
	flag.StringVar(&conf.addr, "address", ":4000", "Server Address")

	flag.StringVar(
		&conf.db.dsn, "dsn",
		os.Getenv("CONTINUITY_DB_DSN"), "Database Dsn",
	)
	flag.IntVar(
		&conf.db.maxIdleConns, "db-max-idle-conns",
		25, "Maximum Idle Db Connections Allowed",
	)
	flag.DurationVar(
		&conf.db.maxIdleTime, "db-max-idle-time",
		15*time.Minute, "Idle Db Connection Time Limit",
	)
	flag.IntVar(
		&conf.db.maxOpenConns, "db-max-open-conns",
		25, "Maximum number of Open Db Connections",
	)
	flag.Parse()
}

func openDb(conf config) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.db.dsn)
	if err != nil {
		return nil, fmt.Errorf("couldn't open db: %w", err)
	}
	db.SetMaxOpenConns(conf.db.maxOpenConns)
	db.SetConnMaxIdleTime(conf.db.maxIdleTime)
	db.SetMaxIdleConns(conf.db.maxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("couldn't ping db: %w", err)
	}
	return db, nil
}
