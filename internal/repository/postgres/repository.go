package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonikq/gravitum_test_task/internal/models"
	"log"
	"strconv"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewStorage(ctx context.Context, uri string, dbPoolWorkers int) (*Repository, error) {
	const source = "repository.NewStorage"
	t1 := time.Now()
	var pool *pgxpool.Pool
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf(models.ErrTraceLayout, source, err)
	}
	config.MaxConns = int32(dbPoolWorkers)

	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf(models.ErrTraceLayout, source, err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf(models.ErrTraceLayout, source, err)
	}
	//err = retrier.DoWithRetries(func() error {
	//	return pool.Ping(ctx)
	//})
	//if err != nil {
	//	return nil, fmt.Errorf(models.ErrTraceLayout, source, err)
	//}

	if err = migrate(pool); err != nil {
		return nil, fmt.Errorf(models.ErrTraceLayout, source, err)
	}

	log.Printf("connection to database took: %v\n", time.Since(t1))

	return &Repository{pool: pool}, nil
}

// Close - closing connection pool to DB.
func (r *Repository) Close() {
	r.pool.Close()
}

// CreateUser - creates a new user.
func (r *Repository) CreateUser(ctx context.Context, body models.UserInfo) (string, error) {
	const source = "repository.CreateUser"
	var userID int64
	err := r.pool.QueryRow(ctx, createUser, body.Username,
		body.FirstName, body.MiddleName, body.LastName, body.Email, body.Gender, body.Age).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return "", models.ErrUsernameIsAlreadyTaken
			}
		}
	}
	return strconv.Itoa(int(userID)), err
}

// GetUser - getting user info by id.
func (r *Repository) GetUser(ctx context.Context, id int64) (*models.UserInfo, error) {
	const source = "repository.GetUser"
	userInfo := new(models.UserInfo)
	if err := r.pool.QueryRow(ctx, getUser, id).
		Scan(&userInfo.ID, &userInfo.Username, &userInfo.FirstName, &userInfo.MiddleName,
			&userInfo.LastName, &userInfo.Email, &userInfo.Gender, &userInfo.Age, &userInfo.EndDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserDoesNotExist
		}

		return nil, fmt.Errorf(models.ErrTraceLayout, source, "error in getting user info: "+err.Error())
	}
	return userInfo, nil
}

// UpdateUser - updating user info by id.
func (r *Repository) UpdateUser(ctx context.Context, body models.UserInfo, id int64) error {
	const source = "repository.UpdateUser"
	_, err := r.pool.Exec(ctx, updateUser, body.Username, body.FirstName,
		body.MiddleName, body.LastName, body.Email, body.Gender, body.Age, id)
	if err != nil {
		return fmt.Errorf(models.ErrTraceLayout, source, "error in updating user info: "+err.Error())
	}
	return nil
}

// DeleteUser - setting end_date for user meta.
func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	const source = "repository.UpdateUser"
	_, err := r.pool.Exec(ctx, deleteUser, id)
	if err != nil {
		return fmt.Errorf(models.ErrTraceLayout, source, "error in deleting user info: "+err.Error())
	}
	return nil
}
