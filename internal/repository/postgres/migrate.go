package postgres

import (
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sonikq/gravitum_test_task/internal/models"
)

//go:embed migration/*.sql
var embedMigrations embed.FS

func migrate(pool *pgxpool.Pool) error {
	const source = "postgres.migrate"
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf(models.ErrTraceLayout, source, "goose.SetDialect: "+err.Error())
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "migration"); err != nil {
		return fmt.Errorf(models.ErrTraceLayout, source, "goose.Up: "+err.Error())
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf(models.ErrTraceLayout, source, err)
	}
	return nil
}
