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
)

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
