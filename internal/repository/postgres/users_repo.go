package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func (u *UserRepository) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (id, user_role, first_name, last_name, position, email, hashed_password, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

	_, err := u.DB.ExecContext(
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

func (u *UserRepository) findUser(query string, args ...any) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user domain.User
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) FindById(userId uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active
		FROM users
		WHERE id = $1		
`
	return u.findUser(query, userId)
}

func (u *UserRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active
		FROM users
		WHERE email = $1		
`
	return u.findUser(query, email)
}

func (u *UserRepository) FindAdmin(admin domain.Role) (*domain.User, error) {
	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active 
		FROM users 
		WHERE user_role = $1 
		LIMIT 1
`
	user, err := u.findUser(query, admin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) Update(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users SET is_active = $1 WHERE id = $2`

	res, err := u.DB.ExecContext(ctx, query, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (u *UserRepository) FindAllByRole(roles ...domain.Role) ([]domain.User, error) {
	if len(roles) == 0 {
		return []domain.User{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, user_role, first_name, last_name, position, email, hashed_password, is_active
		FROM users
		WHERE user_role = ANY($1)
`
	rows, err := u.DB.QueryContext(ctx, query, roles)
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

	if users == nil {
		users = []domain.User{}
	}

	return users, nil
}

func (u *UserRepository) IsActive(userID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT is_active FROM users WHERE id = $1`
	var active bool
	err := u.DB.QueryRowContext(ctx, query, userID).Scan(&active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.ErrUserNotFound
		}
		return false, err
	}

	return active, nil
}
