package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/helf4ch/gocrudl/cmd/api/docs"
	db "github.com/helf4ch/gocrudl/db/migrations"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/config"
	"github.com/helf4ch/gocrudl/internal/handlers"
	"github.com/helf4ch/gocrudl/internal/middlewares"
	"github.com/helf4ch/gocrudl/internal/store"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	httpSwagger "github.com/swaggo/http-swagger"
)

func runMigrations(dbs *sql.DB) error {
	goose.SetBaseFS(db.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(dbs, "."); err != nil {
		return err
	}

	return nil
}

// @title Подписки
// @version 1.0
// @description Микросервис подписок
// @host localhost:8080
// @BasePath /
func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	err := godotenv.Load()
	if err != nil {
		slog.Error("godotenv.Load()", slog.String("err", err.Error()))
		os.Exit(1)
	}

	cfg, err := config.New()
	if err != nil {
		log.Error("config.New()", slog.String("err", err.Error()))
		os.Exit(2)
	}

	str, err := store.New(cfg.DbConn)
	if err != nil {
		log.Error("store.New()", slog.String("err", err.Error()))
		os.Exit(3)
	}

	err = runMigrations(stdlib.OpenDBFromPool(str.Pool))
	if err != nil {
		log.Error("runMigrations", slog.String("err", err.Error()))
		os.Exit(4)
	}

	app := application.New(log, cfg, str)

	app.SetDefaultMiddlewares(
		middlewares.LoggingMiddleware,
		middlewares.RecoverMiddleware,
		middlewares.ErrorMiddleware,
	)

	app.Handle(
		"GET /subs/{id}",
		handlers.Read,
	)

	app.Handle(
		"POST /subs/",
		handlers.Create,
	)

	app.Handle(
		"PUT /subs/{id}",
		handlers.Update,
	)

	app.Handle(
		"DELETE /subs/{id}",
		handlers.Delete,
	)

	app.Handle(
		"GET /subs/list",
		handlers.List,
	)

	app.Handle(
		"GET /subs/total",
		handlers.Total,
	)

	app.RegisterHandler(
		"/swagger/",
		func(
			a application.Application,
			w http.ResponseWriter,
			r *http.Request) error {
			httpSwagger.WrapHandler.ServeHTTP(w, r)
			return nil
		},
	)

	err = app.ListenAndServe()
	if err != nil {
		log.Error("listen and serve", slog.String("err", err.Error()))
		os.Exit(100)
	}
}
