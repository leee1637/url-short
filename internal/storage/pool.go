package storage

import "github.com/jackc/pgx/v5/pgxpool"

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgresPool(pool *pgxpool.Pool) *Postgres {
	return &Postgres{
		pool: pool,
	}
}
