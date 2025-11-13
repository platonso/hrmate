package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/platonso/hrmate/internal/repository"
	"time"
)

type Repository struct {
	Users repository.UserRepository
	Forms repository.FormRepository
	db    *pgxpool.Pool
}

func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	db, err := newPool(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Users: NewUserRepository(db),
		Forms: NewFormRepository(db),
		db:    db,
	}, nil
}

func (r *Repository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

func newPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return pool, nil
}
