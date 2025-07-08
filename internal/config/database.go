package config

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(env *Env) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(env.DatabaseURL)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = 5
	poolCfg.MaxConnLifetime = 30 * 60   // 30 minutes
	poolCfg.HealthCheckPeriod = 10 * 60 // 10 minutes

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	m, err := migrate.New("file://migrations", env.DatabaseURL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		pool.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return pool, nil
}
