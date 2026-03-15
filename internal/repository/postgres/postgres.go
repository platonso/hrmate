package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/platonso/hrmate/internal/repository/postgres/form"
	"github.com/platonso/hrmate/internal/repository/postgres/user"
)

type Repository struct {
	Users user.Repository
	Forms form.Repository
	db    *pgxpool.Pool
}

func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	db, err := NewPool(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Users: *user.NewRepository(db),
		Forms: *form.NewRepository(db),
		db:    db,
	}, nil
}

func (r *Repository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

func NewPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
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
