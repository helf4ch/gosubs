package handlers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	db "github.com/helf4ch/gocrudl/db/sqlc"
	"github.com/helf4ch/gocrudl/internal/application"
	"github.com/helf4ch/gocrudl/internal/dto"
	"github.com/helf4ch/gocrudl/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func List(
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

	priceRaw := r.URL.Query().Get("price")
	priceParsed, err := strconv.Atoi(priceRaw)
	if priceRaw != "" && err != nil {
		return &application.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid price",
			Err:     err,
		}
	}
	var price pgtype.Int4
	if priceRaw != "" {
		price.Int32 = int32(priceParsed)
		price.Valid = true
	}

	serviceNameRaw := r.URL.Query().Get("serviceName")
	var serviceName pgtype.Text
	if serviceNameRaw != "" {
		serviceName.String = serviceNameRaw
		serviceName.Valid = true
	}

	startDateRaw := r.URL.Query().Get("startDate")
	startDate, err := utils.ParseMonthYear(startDateRaw)
	if startDateRaw != "" && err != nil {
		return &application.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid start date",
			Err:     err,
		}
	}

	endDateRaw := r.URL.Query().Get("endDate")
	endDate, err := utils.ParseMonthYear(endDateRaw)
	if endDateRaw != "" && err != nil {
		return &application.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid end date",
			Err:     err,
		}
	}

	subs, err := app.Store.ListSubscriptions(
		r.Context(),
		app.Store.Pool,
		db.ListSubscriptionsParams{
			UserID:      userId,
			Price:       price,
			ServiceName: serviceName,
			StartDate:   startDate,
			EndDate:     endDate,
		},
	)
	if err != nil {
		return err
	}

	var out dto.ListSubscriptionResponse
	for _, sub := range subs {
		out.List = append(out.List, dto.Subscription{
			Id:          sub.ID,
			ServiceName: sub.ServiceName,
			Price:       int(sub.Price),
			UserId:      sub.UserID,
			StartDate:   utils.FormatMonthYear(sub.StartDate),
			EndDate:     utils.FormatMonthYear(sub.EndDate),
			CreatedAt:   sub.CreatedAt.Time.String(),
			UpdatedAt:   utils.FormatMaybeNullTimestamp(sub.UpdatedAt),
		})
	}

	app.ResponseOk(w, out)

	return nil
}
