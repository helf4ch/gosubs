package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/dto"
	"github.com/helf4ch/gocrudl/internal/utils"
)

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
		if err == sql.ErrNoRows {
			return &application.AppError{
				Code:    http.StatusNoContent,
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

	app.ResponseOk(w, out)

	return nil
}
