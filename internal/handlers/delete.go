package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/dto"
)

func Delete(
	app application.Application,
	w http.ResponseWriter,
	r *http.Request,
) error {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return &application.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
			Err:     err,
		}
	}

	count, err := app.Store.DeleteSubscription(r.Context(), app.Store.Pool, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &application.AppError{
				Code:    http.StatusNoContent,
				Message: "no such record",
				Err:     err,
			}
		}
		return err
	}

	out := dto.DeleteSubscriptionResponse{
		Count: int(count),
	}

	err = app.ResponseOk(w, out)
	if err != nil {
		return err
	}

	return nil
}
