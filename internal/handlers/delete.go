package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/dto"
)

// Delete godoc
// @Summary Удалить подписку
// @Description Удалить подписку по id
// @Produce json
// @Param id path uuid.UUID true "id подписки"
// @Success 200 {object} application.AppResponse{body=dto.DeleteSubscriptionResponse}
// @Failure 400 {object} application.AppResponse{error=string}
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

	count, err := app.Store.DeleteSubscription(r.Context(), app.Store.Pool, id)
	if err != nil {
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
