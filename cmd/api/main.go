package main

import (
	"log/slog"
	"os"

	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/config"
	"github.com/helf4ch/gocrudl/internal/handlers"
	"github.com/helf4ch/gocrudl/internal/middlewares"
	"github.com/helf4ch/gocrudl/internal/store"
	"github.com/joho/godotenv"
)

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

	err = app.ListenAndServe()
	if err != nil {
		log.Error("listen and serve", slog.String("err", err.Error()))
		os.Exit(100)
	}
}
