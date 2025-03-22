package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewStorage(ctx context.Context, uri string, dbPoolWorkers int) (*Repository, error) {
	t1 := time.Now()
	var pool *pgxpool.Pool
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(dbPoolWorkers)

	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	if err = migrate(pool); err != nil {
		return nil, err
	}

	log.Printf("connection to database took: %v\n", time.Since(t1))

	return &Repository{pool: pool}, nil
}

func (r *Repository) Close() {
	r.pool.Close()
}
