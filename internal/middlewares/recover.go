package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/helf4ch/gocrudl/internal/application"
)

func RecoverMiddleware(next application.Handler) application.Handler {
	return func(
		app application.Application,
		w http.ResponseWriter,
		r *http.Request,
	) error {
		defer func() {
			err := recover()
			if err != nil {
				app.Log.Error(
					"recover middleware caught panic",
					slog.String("err", fmt.Sprint("%v", err)),
					slog.String("stack", string(debug.Stack())),
				)
				app.Response(w, false, http.StatusInternalServerError, "internal error", nil)
			}
		}()

		return next(app, w, r)
	}
}
