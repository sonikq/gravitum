package repository

import (
	"context"
	"github.com/sonikq/gravitum_test_task/internal/config"
	"github.com/sonikq/gravitum_test_task/internal/models"
	"github.com/sonikq/gravitum_test_task/internal/repository/postgres"
)

type IRepository interface {
	Close()
	CreateUser(ctx context.Context, body models.UserInfo) (string, error)
	GetUser(ctx context.Context, id int64) (*models.UserInfo, error)
	UpdateUser(ctx context.Context, body models.UserInfo, id int64) error
	DeleteUser(ctx context.Context, id int64) error
}

func New(ctx context.Context, cfg config.Config) (IRepository, error) {
	return postgres.NewStorage(ctx, cfg.DatabaseDSN, cfg.DBPoolWorkers)
}
