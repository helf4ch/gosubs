package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/helf4ch/gocrudl/internal/application"
)

func LoggingMiddleware(next application.Handler) application.Handler {
	return func(
		app application.Application,
		w http.ResponseWriter,
		r *http.Request,
	) error {
		requestId := uuid.New().String()

		app.Log = app.Log.With(
			slog.String("requestId", requestId),
			slog.String("pattern", r.URL.RawQuery),
		)

		app.Log.Info("request started")

		ctx := context.WithValue(r.Context(), "requestId", requestId)

		return next(app, w, r.WithContext(ctx))
	}
}
