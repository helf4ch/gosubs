package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	db "github.com/helf4ch/gocrudl/db/sqlc"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/dto"
	"github.com/helf4ch/gocrudl/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"
)

// Create godoc
// @Summary Создать подписку
// @Description Создать подписку
// @Accept json
// @Produce json
// @Param sub body dto.CreateSubscriptionRequest true "тело подписки"
// @Success 200 {object} application.AppResponse{body=dto.CreateSubscriptionResponse}
// @Failure 400 {object} application.AppResponse{error=string}
// @Failure 500 {object} application.AppResponse{error=string}
// @Router /subs/ [post]
func Create(
	app application.Application,
	w http.ResponseWriter,
	r *http.Request,
) error {
	defer r.Body.Close()

	var in dto.CreateSubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&in)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return &application.AppError{
					Code:    http.StatusConflict,
					Message: "record already exists",
					Err:     err,
				}
			}
		}
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
