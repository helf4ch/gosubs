package middlewares

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/helf4ch/gocrudl/internal/application"
)

func ErrorMiddleware(next application.Handler) application.Handler {
	return func(
		app application.Application,
		w http.ResponseWriter,
		r *http.Request,
	) error {
		err := next(app, w, r)

		if err == nil {
			app.Log.Info("handler returned no error")
			return nil
		}

		var appErr *application.AppError

		if errors.As(err, &appErr) {
			app.Response(w, false, appErr.Code, appErr.Message)

			app.Log.Info(
				"handler returned app error",
				slog.String("err", appErr.Err.Error()),
				slog.String("message", appErr.Message),
			)

			return nil
		}

		app.Log.Error(
			"internal error in handler",
			slog.String("err", err.Error()),
		)

		app.ResponseInternalError(w)

		return nil
	}
}
