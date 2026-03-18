package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	db "github.com/helf4ch/gocrudl/db/sqlc"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/dto"
	"github.com/helf4ch/gocrudl/internal/utils"
)

func Create(
	app application.Application,
	w http.ResponseWriter,
	r *http.Request,
) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var in dto.CreateSubscriptionRequest
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

	created, err := app.Store.CreateSubscription(
		r.Context(),
		app.Store.Pool,
		db.CreateSubscriptionParams{
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

	out := dto.CreateSubscriptionResponse{
		Subscription: dto.Subscription{
			Id:          created.ID,
			ServiceName: created.ServiceName,
			Price:       int(created.Price),
			UserId:      created.UserID,
			StartDate:   utils.FormatMonthYear(created.StartDate),
			EndDate:     utils.FormatMonthYear(created.EndDate),
			CreatedAt:   created.CreatedAt.Time.String(),
			UpdatedAt:   utils.FormatMaybeNullTimestamp(created.UpdatedAt),
		},
	}

	err = app.ResponseCreated(w, out)
	if err != nil {
		return err
	}

	return nil
}
