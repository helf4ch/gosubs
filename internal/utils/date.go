package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseMonthYear(dateRaw string) (pgtype.Date, error) {
	if dateRaw == "" {
		return pgtype.Date{Valid: false}, nil
	}

	dateParsed, err := time.Parse("01-2006", dateRaw)
	if err != nil {
		return pgtype.Date{}, err
	}

	return pgtype.Date{
		Time:  dateParsed,
		Valid: true,
	}, nil
}

func FormatMonthYear(date pgtype.Date) string {
	if !date.Valid {
		return ""
	}

	return date.Time.Format("01-2006")
}

func FormatMaybeNullTimestamp(date pgtype.Timestamp) string {
	if !date.Valid {
		return ""
	}

	return date.Time.String()
}
