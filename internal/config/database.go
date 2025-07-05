package config

import (
	"context"

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

	return pool, nil
}
