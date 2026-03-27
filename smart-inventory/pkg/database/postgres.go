package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(conn string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), conn)
}
