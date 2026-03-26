package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/jackc/pgx/v5"
)

// Delete godoc
// @Summary Удалить подписку
// @Description Удалить подписку по id
// @Produce json
// @Param id path uuid.UUID true "id подписки"
// @Success 204 {object} application.AppResponse{}
// @Failure 400 {object} application.AppResponse{error=string}
// @Failure 404 {object} application.AppResponse{error=string}
// @Failure 500 {object} application.AppResponse{error=string}
// @Router /subs/{id} [delete]
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

	_, err = app.Store.DeleteSubscription(r.Context(), app.Store.Pool, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &application.AppError{
				Code:    http.StatusNotFound,
				Message: "no such record",
				Err:     err,
			}
		}
		return err
	}

	err = app.Response(w, true, http.StatusNoContent, "", nil)
	if err != nil {
		return err
	}

	return nil
}
