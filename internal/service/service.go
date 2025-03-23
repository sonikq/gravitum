package service

import (
	"context"
	"github.com/sonikq/gravitum_test_task/internal/models"
	"github.com/sonikq/gravitum_test_task/internal/repository"
	"github.com/sonikq/gravitum_test_task/internal/service/user_management"
)

type IUserManagementService interface {
	CreateUser(ctx context.Context, request models.UserInfo) (string, error)
	GetUser(ctx context.Context, id int64) (*models.UserInfo, error)
	UpdateUser(ctx context.Context, request models.UserInfo) error
	DeleteUser(ctx context.Context, id int64) error
}

type Service struct {
	IUserManagementService
}

func New(repo repository.IRepository) *Service {
	return &Service{
		IUserManagementService: user_management.NewService(repo),
	}
}
