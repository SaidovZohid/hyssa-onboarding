package repository

import (
	"context"

	"learning/hyssa-learn/internal/core/repository/psql"
	"learning/hyssa-learn/internal/core/repository/psql/sqlc"
)

type Store interface {
	sqlc.Querier
}

func New(ctx context.Context, dsn string) Store {
	return psql.NewStore(ctx, dsn)
}
