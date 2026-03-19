package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	db "github.com/helf4ch/gocrudl/db/sqlc"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/dto"
	"github.com/helf4ch/gocrudl/internal/utils"
	"github.com/jackc/pgx/v5"
)

// Update godoc
// @Summary Редактировать подписку
// @Description Редактировать подписку по id
// @Accept json
// @Produce json
// @Param id path uuid.UUID true "id подписки"
// @Param sub body dto.UpdateSubscriptionRequest true "новое тело подписки"
// @Success 200 {object} application.AppResponse{body=dto.UpdateSubscriptionResponse}
// @Failure 400 {object} application.AppResponse{error=string}
// @Failure 404 {object} application.AppResponse{error=string}
// @Failure 500 {object} application.AppResponse{error=string}
// @Router /subs/{id} [put]
func Update(
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

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var in dto.UpdateSubscriptionRequest
	err = json.Unmarshal(body, &in)
	if err != nil {
		return &application.AppError{
			Code:    http.StatusBadRequest,
			Message: "bad formed input data",
			Err:     err,
		}
	}

	inStartDate, err := utils.ParseMonthYear(in.StartDate)
	if err != nil {
		return err
	}

	inEndDate, err := utils.ParseMonthYear(in.EndDate)
	if err != nil {
		return err
	}

	sub, err := app.Store.UpdateSubscription(
		r.Context(),
		app.Store.Pool,
		db.UpdateSubscriptionParams{
			ID:          id,
			ServiceName: in.ServiceName,
			Price:       int32(in.Price),
			UserID:      in.UserId,
			StartDate:   inStartDate,
			EndDate:     inEndDate,
		},
	)
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

	out := dto.UpdateSubscriptionResponse{
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
