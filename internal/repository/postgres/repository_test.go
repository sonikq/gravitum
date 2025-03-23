package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sonikq/gravitum_test_task/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupPostgresContainer sets up a PostgreSQL container for testing
func setupPostgresContainer(ctx context.Context) (testcontainers.Container, string, error) {
	// Define PostgreSQL container configuration
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	// Start the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start container: %w", err)
	}

	// Get the container's host and port
	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get container port: %w", err)
	}

	// Construct the connection string
	connString := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())

	return container, connString, nil
}

func TestRepository_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Setup PostgreSQL container
	container, connString, err := setupPostgresContainer(ctx)
	require.NoError(t, err)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}()

	// Create repository
	repo, err := NewStorage(ctx, connString, 5)
	require.NoError(t, err)
	defer repo.Close()

	// Run the tests
	t.Run("CreateAndGetUser", func(t *testing.T) {
		testCreateAndGetUser(ctx, t, repo)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		testUpdateUser(ctx, t, repo)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		testDeleteUser(ctx, t, repo)
	})

	t.Run("CreateDuplicateUser", func(t *testing.T) {
		testCreateDuplicateUser(ctx, t, repo)
	})

	t.Run("GetNonExistentUser", func(t *testing.T) {
		testGetNonExistentUser(ctx, t, repo)
	})
}

func testCreateAndGetUser(ctx context.Context, t *testing.T, repo *Repository) {
	// Create a test user
	user := models.UserInfo{
		Username:   "testuser",
		FirstName:  "Test",
		MiddleName: "Middle",
		LastName:   "User",
		Email:      "test@example.com",
		Gender:     "M",
		Age:        30,
	}

	// Create the user
	result, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Get the user by ID (assuming ID is 1 for the first user)
	retrievedUser, err := repo.GetUser(ctx, 1)
	require.NoError(t, err)

	// Verify user data
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.FirstName, retrievedUser.FirstName)
	assert.Equal(t, user.MiddleName, retrievedUser.MiddleName)
	assert.Equal(t, user.LastName, retrievedUser.LastName)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.Gender, retrievedUser.Gender)
	assert.Equal(t, user.Age, retrievedUser.Age)
	assert.Nil(t, retrievedUser.EndDate)
}

func testUpdateUser(ctx context.Context, t *testing.T, repo *Repository) {
	// Create a test user first
	user := models.UserInfo{
		Username:   "updateuser",
		FirstName:  "Update",
		MiddleName: "Middle",
		LastName:   "User",
		Email:      "update@example.com",
		Gender:     "F",
		Age:        25,
	}

	// Create the user
	result, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Get the user to find the ID
	retrievedUser, err := repo.GetUser(ctx, 2) // Assuming this is the second user
	require.NoError(t, err)

	// Update user data
	updatedUser := models.UserInfo{
		Username:   "updateduser",
		FirstName:  "Updated",
		MiddleName: "NewMiddle",
		LastName:   "UserUpdated",
		Email:      "updated@example.com",
		Gender:     "M",
		Age:        35,
	}

	// Update the user
	err = repo.UpdateUser(ctx, updatedUser, retrievedUser.ID)
	require.NoError(t, err)

	// Get the updated user
	retrievedUpdatedUser, err := repo.GetUser(ctx, retrievedUser.ID)
	require.NoError(t, err)

	// Verify updated user data
	assert.Equal(t, updatedUser.Username, retrievedUpdatedUser.Username)
	assert.Equal(t, updatedUser.FirstName, retrievedUpdatedUser.FirstName)
	assert.Equal(t, updatedUser.MiddleName, retrievedUpdatedUser.MiddleName)
	assert.Equal(t, updatedUser.LastName, retrievedUpdatedUser.LastName)
	assert.Equal(t, updatedUser.Email, retrievedUpdatedUser.Email)
	assert.Equal(t, updatedUser.Gender, retrievedUpdatedUser.Gender)
	assert.Equal(t, updatedUser.Age, retrievedUpdatedUser.Age)
}

func testDeleteUser(ctx context.Context, t *testing.T, repo *Repository) {
	// Create a test user first
	user := models.UserInfo{
		Username:   "deleteuser",
		FirstName:  "Delete",
		MiddleName: "Middle",
		LastName:   "User",
		Email:      "delete@example.com",
		Gender:     "M",
		Age:        40,
	}

	// Create the user
	result, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Get the user to find the ID
	retrievedUser, err := repo.GetUser(ctx, 3) // Assuming this is the third user
	require.NoError(t, err)

	// Delete the user
	err = repo.DeleteUser(ctx, retrievedUser.ID)
	require.NoError(t, err)

	// Get the deleted user
	deletedUser, err := repo.GetUser(ctx, retrievedUser.ID)
	require.NoError(t, err)

	// Verify the user has an end date set (soft delete)
	assert.NotNil(t, deletedUser.EndDate)
	assert.WithinDuration(t, time.Now(), *deletedUser.EndDate, 5*time.Second)
}

func testCreateDuplicateUser(ctx context.Context, t *testing.T, repo *Repository) {
	// Create a test user
	user := models.UserInfo{
		Username:   "duplicateuser",
		FirstName:  "Duplicate",
		MiddleName: "Middle",
		LastName:   "User",
		Email:      "duplicate@example.com",
		Gender:     "F",
		Age:        28,
	}

	// Create the user
	result, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Try to create the same user again
	_, err = repo.CreateUser(ctx, user)
	assert.ErrorIs(t, err, models.ErrUsernameIsAlreadyTaken)
}

func testGetNonExistentUser(ctx context.Context, t *testing.T, repo *Repository) {
	// Try to get a user with a non-existent ID
	_, err := repo.GetUser(ctx, 9999)
	assert.ErrorIs(t, err, models.ErrUserDoesNotExist)
}
