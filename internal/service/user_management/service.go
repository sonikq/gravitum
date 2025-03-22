package user_management

import "github.com/sonikq/gravitum_test_task/internal/repository"

type Service struct {
	repository repository.IRepository
}

func NewService(repo repository.IRepository) *Service {
	return &Service{
		repository: repo,
	}
}
