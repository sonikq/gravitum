package user_management

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sonikq/gravitum_test_task/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepository is a mock implementation of the repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Close() {
	return
}

func (m *MockRepository) CreateUser(ctx context.Context, user models.UserInfo) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) GetUser(ctx context.Context, id int64) (*models.UserInfo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserInfo), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, user models.UserInfo, id int64) error {
	args := m.Called(ctx, user, id)
	return args.Error(0)
}

func (m *MockRepository) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to create a valid user for testing
func createValidUser() models.UserInfo {
	return models.UserInfo{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Gender:    "M",
		Age:       20,
		EndDate:   nil,
	}
}

// TestCreateUser tests the CreateUser method
func TestCreateUser(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	service := &Service{repository: mockRepo}
	ctx := context.Background()

	t.Run("Success - Valid user creation", func(t *testing.T) {
		// Arrange
		user := createValidUser()
		expectedID := "user123"
		mockRepo.On("CreateUser", ctx, user).Return(expectedID, nil).Once()

		// Act
		id, err := service.CreateUser(ctx, user)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedID, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Invalid user data", func(t *testing.T) {
		// Arrange
		invalidUser := models.UserInfo{
			// Missing required fields
		}
		// We expect Validate to fail, so repository won't be called

		// Act
		id, err := service.CreateUser(ctx, invalidUser)

		// Assert
		require.Error(t, err)
		assert.Empty(t, id)
		// No need to verify mockRepo as it shouldn't be called
	})

	t.Run("Failure - Repository error", func(t *testing.T) {
		// Arrange
		user := createValidUser()
		expectedError := errors.New("database error")
		mockRepo.On("CreateUser", ctx, user).Return("", expectedError).Once()

		// Act
		id, err := service.CreateUser(ctx, user)

		// Assert
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Edge case - Empty but valid user", func(t *testing.T) {
		// This test depends on what Validate() considers valid
		// For this example, we'll assume minimal valid data
		minimalUser := models.UserInfo{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
			Gender:    "F",
			Age:       24,
		}
		expectedID := "user456"
		mockRepo.On("CreateUser", ctx, minimalUser).Return(expectedID, nil).Once()

		// Act
		id, err := service.CreateUser(ctx, minimalUser)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedID, id)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetUser tests the GetUser method
func TestGetUser(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	service := &Service{repository: mockRepo}
	ctx := context.Background()

	t.Run("Success - Get active user", func(t *testing.T) {
		// Arrange
		userID := int64(1)
		expectedUser := createValidUser()
		mockRepo.On("GetUser", ctx, userID).Return(&expectedUser, nil).Once()

		// Act
		user, err := service.GetUser(ctx, userID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, &expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - User not found", func(t *testing.T) {
		// Arrange
		userID := int64(999)
		expectedError := errors.New("user not found")
		mockRepo.On("GetUser", ctx, userID).Return(nil, expectedError).Once()

		// Act
		user, err := service.GetUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - User is gone (has end date)", func(t *testing.T) {
		// Arrange
		userID := int64(2)
		endDate := time.Now()
		userWithEndDate := createValidUser()
		userWithEndDate.EndDate = &endDate
		mockRepo.On("GetUser", ctx, userID).Return(&userWithEndDate, nil).Once()

		// Act
		user, err := service.GetUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Equal(t, models.ErrUserIsGone, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Repository error", func(t *testing.T) {
		// Arrange
		userID := int64(3)
		expectedError := errors.New("database error")
		mockRepo.On("GetUser", ctx, userID).Return(nil, expectedError).Once()

		// Act
		user, err := service.GetUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Edge case - Zero ID", func(t *testing.T) {
		// Arrange
		userID := int64(0)
		// Behavior depends on repository implementation
		// For this test, we'll assume it returns an error
		expectedError := errors.New("invalid user ID")
		mockRepo.On("GetUser", ctx, userID).Return(nil, expectedError).Once()

		// Act
		user, err := service.GetUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateUser tests the UpdateUser method
func TestUpdateUser(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	service := &Service{repository: mockRepo}
	ctx := context.Background()

	t.Run("Success - Update existing user", func(t *testing.T) {
		// Arrange
		user := createValidUser()
		// First GetUser call to check if user exists
		mockRepo.On("GetUser", ctx, user.ID).Return(&user, nil).Once()
		// Then UpdateUser call
		mockRepo.On("UpdateUser", ctx, user, user.ID).Return(nil).Once()

		// Act
		err := service.UpdateUser(ctx, user)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - User not found", func(t *testing.T) {
		// Arrange
		user := createValidUser()
		user.ID = 999
		expectedError := errors.New("user not found")
		mockRepo.On("GetUser", ctx, user.ID).Return(nil, expectedError).Once()

		// Act
		err := service.UpdateUser(ctx, user)

		// Assert
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - User is gone", func(t *testing.T) {
		// Arrange
		user := createValidUser()
		user.ID = 2
		endDate := time.Now()
		userWithEndDate := createValidUser()
		userWithEndDate.EndDate = &endDate
		mockRepo.On("GetUser", ctx, user.ID).Return(&userWithEndDate, nil).Once()

		// Act
		err := service.UpdateUser(ctx, user)

		// Assert
		require.Error(t, err)
		assert.Equal(t, models.ErrUserIsGone, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Invalid update data", func(t *testing.T) {
		// Arrange
		validUser := createValidUser()
		mockRepo.On("GetUser", ctx, validUser.ID).Return(&validUser, nil).Once()

		// Now create an invalid version for the update
		invalidUser := validUser
		invalidUser.Email = "not-an-email" // Assuming this fails validation

		// Act
		err := service.UpdateUser(ctx, invalidUser)

		// Assert
		require.Error(t, err)
		// No need to verify UpdateUser call as it shouldn't happen
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Repository update error", func(t *testing.T) {
		// Arrange
		user := createValidUser()
		expectedError := errors.New("database error")
		mockRepo.On("GetUser", ctx, user.ID).Return(&user, nil).Once()
		mockRepo.On("UpdateUser", ctx, user, user.ID).Return(expectedError).Once()

		// Act
		err := service.UpdateUser(ctx, user)

		// Assert
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Edge case - Update with minimal changes", func(t *testing.T) {
		// Arrange
		originalUser := createValidUser()
		mockRepo.On("GetUser", ctx, originalUser.ID).Return(&originalUser, nil).Once()

		// Minimal update - just change one field
		updatedUser := originalUser
		updatedUser.LastName = "Smith"
		mockRepo.On("UpdateUser", ctx, updatedUser, updatedUser.ID).Return(nil).Once()

		// Act
		err := service.UpdateUser(ctx, updatedUser)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestDeleteUser tests the DeleteUser method
func TestDeleteUser(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	service := &Service{repository: mockRepo}
	ctx := context.Background()

	t.Run("Success - Delete existing user", func(t *testing.T) {
		// Arrange
		userID := int64(1)
		user := createValidUser()
		mockRepo.On("GetUser", ctx, userID).Return(&user, nil).Once()
		mockRepo.On("DeleteUser", ctx, userID).Return(nil).Once()

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - User not found", func(t *testing.T) {
		// Arrange
		userID := int64(999)
		expectedError := errors.New("user not found")
		mockRepo.On("GetUser", ctx, userID).Return(nil, expectedError).Once()

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - User already deleted", func(t *testing.T) {
		// Arrange
		userID := int64(2)
		endDate := time.Now()
		userWithEndDate := createValidUser()
		userWithEndDate.EndDate = &endDate
		mockRepo.On("GetUser", ctx, userID).Return(&userWithEndDate, nil).Once()

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Equal(t, models.ErrDeleteDeletedUser, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Repository delete error", func(t *testing.T) {
		// Arrange
		userID := int64(3)
		user := createValidUser()
		expectedError := errors.New("database error")
		mockRepo.On("GetUser", ctx, userID).Return(&user, nil).Once()
		mockRepo.On("DeleteUser", ctx, userID).Return(expectedError).Once()

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Edge case - Zero ID", func(t *testing.T) {
		// Arrange
		userID := int64(0)
		// Behavior depends on repository implementation
		// For this test, we'll assume it returns an error
		expectedError := errors.New("invalid user ID")
		mockRepo.On("GetUser", ctx, userID).Return(nil, expectedError).Once()

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestWithCanceledContext tests behavior with canceled context
func TestWithCanceledContext(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	service := &Service{repository: mockRepo}

	// Create canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	t.Run("CreateUser with canceled context", func(t *testing.T) {
		// Arrange
		user := createValidUser()
		expectedError := context.Canceled
		mockRepo.On("CreateUser", ctx, user).Return("", expectedError).Once()

		// Act
		id, err := service.CreateUser(ctx, user)

		// Assert
		require.Error(t, err)
		assert.Empty(t, id)
		mockRepo.AssertExpectations(t)
	})

	// Similar tests for other methods...
}

// TestWithDeadlineExceededContext tests behavior with deadline exceeded context
func TestWithDeadlineExceededContext(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	service := &Service{repository: mockRepo}

	// Create context with immediate deadline
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond) // Ensure deadline is exceeded

	t.Run("GetUser with deadline exceeded", func(t *testing.T) {
		// Arrange
		userID := int64(1)
		expectedError := context.DeadlineExceeded
		mockRepo.On("GetUser", ctx, userID).Return(nil, expectedError).Once()

		// Act
		user, err := service.GetUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	// Similar tests for other methods...
}
