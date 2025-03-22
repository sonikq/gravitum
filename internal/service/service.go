package service

import (
	"context"
	"github.com/sonikq/gravitum_test_task/internal/repository"
	"github.com/sonikq/gravitum_test_task/internal/service/user_management"
)

type IUserManagementService interface {
	Register(ctx context.Context) error
}

type Service struct {
	IUserManagementService
}

func New(repo repository.IRepository) *Service {
	return &Service{
		IUserManagementService: user_management.NewService(repo),
	}
}
