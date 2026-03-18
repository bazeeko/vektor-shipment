package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const defaultMigrationsPath = "migrations"

type Params struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

func NewConnectWithMigration(ctx context.Context, params Params) (*pgxpool.Pool, error) {
	pgConnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		params.User,
		params.Password,
		params.Host,
		params.Port,
		params.DBName,
	)

	pool, err := pgxpool.New(ctx, pgConnString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pool.Ping: %w", err)
	}

	if err = migrate(ctx, pool); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return pool, nil
}

func migrate(ctx context.Context, pool *pgxpool.Pool) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %s", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	if err := goose.UpContext(ctx, db, defaultMigrationsPath); err != nil {
		return fmt.Errorf("goose.UpContext: %w", err)
	}

	return nil
}
