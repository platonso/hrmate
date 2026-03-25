package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	"github.com/platonso/hrmate/internal/repository/postgres/user/entity"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	rec := entity.ToUserRecord(*user)
	query := `
		INSERT INTO users (id, user_role, first_name, last_name, position, email, hashed_password, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`
	_, err := r.DB.Exec(
		ctx,
		query,
		rec.ID,
		rec.Role,
		rec.FirstName,
		rec.LastName,
		rec.Position,
		rec.Email,
		rec.HashedPassword,
		rec.IsActive,
	)
	return err
}

func (r *Repository) FindByUserID(ctx context.Context, userId uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active
		FROM users
		WHERE id = $1		
`
	rec, err := r.findUser(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	user := entity.ToDomainUser(rec)
	return &user, nil
}

func (r *Repository) FindByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]domain.User, error) {
	if len(userIDs) == 0 {
		return []domain.User{}, nil
	}

	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active
		FROM users
		WHERE id = ANY($1)
	`

	rows, err := r.DB.Query(ctx, query, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]entity.UserRecord, 0, len(userIDs))
	for rows.Next() {
		var rec entity.UserRecord
		err := rows.Scan(
			&rec.ID,
			&rec.Role,
			&rec.FirstName,
			&rec.LastName,
			&rec.Position,
			&rec.Email,
			&rec.HashedPassword,
			&rec.IsActive,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entity.ToDomainUsers(records), nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active
		FROM users
		WHERE email = $1		
`
	rec, err := r.findUser(ctx, query, email)
	if err != nil {
		return nil, err
	}
	user := entity.ToDomainUser(rec)
	return &user, nil
}

func (r *Repository) Update(ctx context.Context, user *domain.User) error {
	rec := entity.ToUserRecord(*user)
	query := `
        UPDATE users SET
            user_role = $1,
            first_name = $2,
            last_name = $3,
            position = $4,
            email = $5,
            hashed_password = $6,
            is_active = $7
        WHERE id = $8
    `

	tag, err := r.DB.Exec(ctx, query,
		rec.Role,
		rec.FirstName,
		rec.LastName,
		rec.Position,
		rec.Email,
		rec.HashedPassword,
		rec.IsActive,
		rec.ID,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *Repository) FindByRole(ctx context.Context, roles ...domain.Role) ([]domain.User, error) {
	if len(roles) == 0 {
		return []domain.User{}, nil
	}

	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active
		FROM users
		WHERE user_role = ANY($1)
`
	rows, err := r.DB.Query(ctx, query, roles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]entity.UserRecord, 0)
	for rows.Next() {
		var rec entity.UserRecord
		err := rows.Scan(
			&rec.ID,
			&rec.Role,
			&rec.FirstName,
			&rec.LastName,
			&rec.Position,
			&rec.Email,
			&rec.HashedPassword,
			&rec.IsActive,
		)
		if err != nil {
			return nil, err
		}

		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	users := entity.ToDomainUsers(records)
	return users, nil
}

func (r *Repository) IsActive(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT is_active FROM users WHERE id = $1`

	var active bool
	err := r.DB.QueryRow(ctx, query, userID).Scan(&active)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, errs.ErrUserNotFound
		}
		return false, err
	}

	return active, nil
}

func (r *Repository) findUser(ctx context.Context, query string, args ...any) (entity.UserRecord, error) {
	var rec entity.UserRecord
	err := r.DB.QueryRow(ctx, query, args...).Scan(
		&rec.ID,
		&rec.Role,
		&rec.FirstName,
		&rec.LastName,
		&rec.Position,
		&rec.Email,
		&rec.HashedPassword,
		&rec.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserRecord{}, errs.ErrUserNotFound
		}
		return entity.UserRecord{}, err
	}
	return rec, nil
}
