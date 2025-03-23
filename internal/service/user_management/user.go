package user_management

import (
	"context"
	"github.com/sonikq/gravitum_test_task/internal/models"
)

// CreateUser - validating request body and creating user in DB.
func (s *Service) CreateUser(ctx context.Context, request models.UserInfo) (string, error) {
	if err := request.Validate(); err != nil {
		return "", err
	}

	id, err := s.repository.CreateUser(ctx, request)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetUser - getting info about user by id.
func (s *Service) GetUser(ctx context.Context, id int64) (*models.UserInfo, error) {
	userInfo, err := s.repository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	if userInfo.EndDate != nil {
		return nil, models.ErrUserIsGone
	}

	return userInfo, nil
}

// UpdateUser - updating user info by id.
func (s *Service) UpdateUser(ctx context.Context, request models.UserInfo) error {
	_, err := s.GetUser(ctx, request.ID)
	if err != nil {
		return err
	}

	if err = request.Validate(); err != nil {
		return err
	}

	return s.repository.UpdateUser(ctx, request, request.ID)
}

// DeleteUser - deleting user by id.
func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	userInfo, err := s.repository.GetUser(ctx, id)
	if err != nil {
		return err
	}

	if userInfo.EndDate != nil {
		return models.ErrDeleteDeletedUser
	}

	return s.repository.DeleteUser(ctx, id)
}
