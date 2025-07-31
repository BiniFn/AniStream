package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/coeeter/aniways/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(env *config.Env, log *slog.Logger) (*pgxpool.Pool, error) {
	log.Info("initialising database connection")

	cfg, err := pgxpool.ParseConfig(env.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	cfg.MaxConns = 10
	cfg.MinConns = 2

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("open pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db ping: %w", err)
	}
	log.Info("database ping OK")

	m, err := migrate.New("file://migrations", env.DatabaseURL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("migrate init: %w", err)
	}
	switch err := m.Up(); err {
	case nil:
		log.Info("migrations applied")
	case migrate.ErrNoChange:
		log.Info("migrations already up‑to‑date")
	default:
		pool.Close()
		return nil, fmt.Errorf("migrate up: %w", err)
	}

	log.Info("database ready",
		"max_conns", cfg.MaxConns,
		"min_conns", cfg.MinConns,
	)
	return pool, nil
}
