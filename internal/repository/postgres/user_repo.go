package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/platonso/hrmate/internal/domain"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (u *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, user_role, first_name, last_name, position, email, hashed_password, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`
	_, err := u.DB.Exec(
		ctx,
		query,
		user.ID,
		user.Role,
		user.FirstName,
		user.LastName,
		user.Position,
		user.Email,
		user.HashedPassword,
		user.IsActive,
	)
	return err
}

func (u *UserRepository) findUser(ctx context.Context, query string, args ...any) (*domain.User, error) {
	var user domain.User
	err := u.DB.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Role,
		&user.FirstName,
		&user.LastName,
		&user.Position,
		&user.Email,
		&user.HashedPassword,
		&user.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) FindByUserId(ctx context.Context, userId uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active
		FROM users
		WHERE id = $1		
`
	return u.findUser(ctx, query, userId)
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active
		FROM users
		WHERE email = $1		
`
	return u.findUser(ctx, query, email)
}

func (u *UserRepository) FindAdmin(ctx context.Context) (*domain.User, error) {
	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active 
		FROM users 
		WHERE user_role = $1 
		LIMIT 1
`
	user, err := u.findUser(ctx, query, domain.RoleAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET is_active = $1 WHERE id = $2`

	tag, err := u.DB.Exec(ctx, query, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (u *UserRepository) FindAllByRole(ctx context.Context, roles ...domain.Role) ([]domain.User, error) {
	if len(roles) == 0 {
		return []domain.User{}, nil
	}

	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active
		FROM users
		WHERE user_role = ANY($1)
`
	rows, err := u.DB.Query(ctx, query, roles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Role,
			&user.FirstName,
			&user.LastName,
			&user.Position,
			&user.Email,
			&user.HashedPassword,
			&user.IsActive,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserRepository) IsActive(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT is_active FROM users WHERE id = $1`

	var active bool
	err := u.DB.QueryRow(ctx, query, userID).Scan(&active)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, domain.ErrUserNotFound
		}
		return false, err
	}

	return active, nil
}
