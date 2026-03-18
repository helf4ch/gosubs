package handlers

import (
	"net/http"

	"github.com/google/uuid"
	db "github.com/helf4ch/gocrudl/db/sqlc"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/dto"
	"github.com/helf4ch/gocrudl/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func Total(
	app application.Application,
	w http.ResponseWriter,
	r *http.Request,
) error {
	userIdRaw := r.URL.Query().Get("userId")
	userIdParsed, err := uuid.Parse(userIdRaw)
	if userIdRaw != "" && err != nil {
		return &application.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid userId",
			Err:     nil,
		}
	}
	var userId pgtype.UUID
	if userIdRaw != "" {
		userId.Bytes = userIdParsed
		userId.Valid = true
	}

	serviceNameRaw := r.URL.Query().Get("serviceName")
	var serviceName pgtype.Text
	if serviceNameRaw != "" {
		serviceName.String = serviceNameRaw
		serviceName.Valid = true
	}

	startDateRaw := r.URL.Query().Get("startDate")
	if startDateRaw == "" {
		return &application.AppError{
			Code:    http.StatusBadRequest,
			Message: "empty start date",
			Err:     nil,
		}
	}
	startDate, err := utils.ParseMonthYear(startDateRaw)
	if err != nil {
		return &application.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid start date",
			Err:     err,
		}
	}

	endDateRaw := r.URL.Query().Get("endDate")
	if endDateRaw == "" {
		return &application.AppError{
			Code:    http.StatusBadRequest,
			Message: "empty end date",
			Err:     nil,
		}
	}
	endDate, err := utils.ParseMonthYear(endDateRaw)
	if endDateRaw != "" && err != nil {
		return &application.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid end date",
			Err:     err,
		}
	}

	cost, err := app.Store.TotalSubscriptionsCost(
		r.Context(),
		app.Store.Pool,
		db.TotalSubscriptionsCostParams{
			UserID:      userId,
			ServiceName: serviceName,
			StartDate:   startDate,
			EndDate:     endDate,
		},
	)
	if err != nil {
		return err
	}

	out := dto.TotalSubscriptionResponse{
		Cost: int(cost),
	}

	app.ResponseOk(w, out)

	return nil
}
