package config

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(env *Env) (*pgxpool.Pool, error) {
	log.Printf("🚀 Initializing database connection…")

	// open the pool
	poolCfg, err := pgxpool.ParseConfig(env.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DATABASE_URL: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection pool: %w", err)
	}

	// verify connectivity
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	log.Printf("✅ Connected to database")

	// apply migrations
	m, err := migrate.New("file://migrations", env.DatabaseURL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("migration init error: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		pool.Close()
		return nil, fmt.Errorf("migration apply error: %w", err)
	}
	log.Printf("✨ Migrations applied")

	log.Printf("🎉 Database ready (maxConns=%d)", poolCfg.MaxConns)
	return pool, nil
}
