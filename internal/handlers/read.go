package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/dto"
	"github.com/helf4ch/gocrudl/internal/utils"
	"github.com/jackc/pgx/v5"
)

// Read godoc
// @Summary Получить подписку
// @Description Получить подписку по id
// @Produce json
// @Param id path uuid.UUID true "id подписки"
// @Success 200 {object} application.AppResponse{body=dto.ReadSubscriptionResponse}
// @Failure 400 {object} application.AppResponse{error=string}
// @Failure 404 {object} application.AppResponse{error=string}
// @Failure 500 {object} application.AppResponse{error=string}
// @Router /subs/{id} [get]
func Read(
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

	sub, err := app.Store.ReadSubscription(r.Context(), app.Store.Pool, id)
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

	out := dto.ReadSubscriptionResponse{
		Subscription: dto.Subscription{
			Id:          sub.ID,
			ServiceName: sub.ServiceName,
			Price:       int(sub.Price),
			UserId:      sub.UserID,
			StartDate:   utils.FormatMonthYear(sub.StartDate),
			EndDate:     utils.FormatMonthYear(sub.EndDate),
			CreatedAt:   sub.CreatedAt.Time.String(),
			UpdatedAt:   utils.FormatMaybeNullTimestamp(sub.UpdatedAt),
		},
	}

	err = app.ResponseOk(w, out)
	if err != nil {
		return err
	}

	return nil
}
