package store

import (
	"context"

	db "github.com/helf4ch/gocrudl/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
	db.Querier
}

func New(dbUrl string) (*Store, error) {
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	return &Store{
		pool,
		db.New(),
	}, nil
}
