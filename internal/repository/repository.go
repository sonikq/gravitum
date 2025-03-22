package repository

import (
	"context"
	"github.com/sonikq/gravitum_test_task/internal/config"
	"github.com/sonikq/gravitum_test_task/internal/repository/postgres"
)

type IRepository interface {
	Close()
}

func New(ctx context.Context, cfg config.Config) (IRepository, error) {
	return postgres.NewStorage(ctx, cfg.DatabaseDSN, cfg.DBPoolWorkers)
}
