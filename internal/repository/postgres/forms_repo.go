package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	"time"
)

type FormRepository struct {
	DB *sql.DB
}

func (f *FormRepository) Create(form *domain.Form) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO forms (id, user_id, title, description, start_date, end_date, created_at, approved_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`
	_, err := f.DB.ExecContext(
		ctx,
		query,
		form.ID,
		form.UserID,
		form.Title,
		form.Description,
		form.StartDate,
		form.EndDate,
		form.CreatedAt,
		form.ApprovedAt,
		form.Status,
	)
	return err
}

func (f *FormRepository) FindByID(formId uuid.UUID) (*domain.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, title, description, start_date, end_date, created_at, approved_at, status
		FROM forms
		WHERE id = $1
	`

	var form domain.Form
	err := f.DB.QueryRowContext(ctx, query, formId).Scan(
		&form.ID,
		&form.UserID,
		&form.Title,
		&form.Description,
		&form.StartDate,
		&form.EndDate,
		&form.CreatedAt,
		&form.ApprovedAt,
		&form.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrFormNotFound
		}
		return nil, err
	}

	return &form, nil
}

func (f *FormRepository) FindByUserIDWithUser(userID uuid.UUID) ([]domain.UserWithForms, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT 
		    u.id, u.first_name, u.last_name, u.position, u.email,
		    f.id, f.user_id, f.title, f.description, f.start_date, 
		    f.end_date, f.created_at, f.approved_at, f.status
		FROM forms f 
		JOIN users u ON u.id = f.user_id
		WHERE u.id = $1
`

	rows, err := f.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user domain.User
	var forms []domain.Form

	for rows.Next() {
		var form domain.Form
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Position,
			&user.Email,
			&form.ID,
			&form.UserID,
			&form.Title,
			&form.Description,
			&form.StartDate,
			&form.EndDate,
			&form.CreatedAt,
			&form.ApprovedAt,
			&form.Status,
		)
		if err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if user.ID == uuid.Nil {
		return []domain.UserWithForms{}, nil
	}

	result := []domain.UserWithForms{
		{
			User:  user,
			Forms: forms,
		},
	}

	return result, nil

}

func (f *FormRepository) FindByIDWithUser(formId uuid.UUID) (*domain.UserWithForm, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT 
		    u.id, u.first_name, u.last_name, u.position, u.email,
		    f.id, f.user_id, f.title, f.description, f.start_date, 
		    f.end_date, f.created_at, f.approved_at, f.status
		FROM forms f 
		JOIN users u ON u.id = f.user_id
		WHERE f.id = $1
`
	var res domain.UserWithForm

	err := f.DB.QueryRowContext(ctx, query, formId).Scan(
		&res.User.ID,
		&res.User.FirstName,
		&res.User.LastName,
		&res.User.Position,
		&res.User.Email,
		&res.Form.ID,
		&res.Form.UserID,
		&res.Form.Title,
		&res.Form.Description,
		&res.Form.StartDate,
		&res.Form.EndDate,
		&res.Form.CreatedAt,
		&res.Form.ApprovedAt,
		&res.Form.Status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrFormNotFound
		}
		return nil, err
	}

	return &res, nil
}

func (f *FormRepository) FindAllWithUsers() ([]domain.UserWithForms, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT 
			u.id, u.first_name, u.last_name, u.position, u.email,
		    f.id, f.user_id, f.title, f.description, f.start_date, 
		    f.end_date, f.created_at, f.approved_at, f.status
		FROM users u
		JOIN forms f ON u.id = f.user_id
`
	rows, err := f.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersMap = make(map[uuid.UUID]*domain.UserWithForms)

	for rows.Next() {
		var user domain.User
		var form domain.Form

		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Position, &user.Email,
			&form.ID, &form.UserID, &form.Title, &form.Description, &form.StartDate, &form.EndDate,
			&form.CreatedAt, &form.ApprovedAt, &form.Status,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := usersMap[user.ID]; !ok {
			usersMap[user.ID] = &domain.UserWithForms{
				User:  user,
				Forms: []domain.Form{},
			}
		}

		usersMap[user.ID].Forms = append(usersMap[user.ID].Forms, form)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var result []domain.UserWithForms
	for _, u := range usersMap {
		result = append(result, *u)
	}

	return result, nil
}

func (f *FormRepository) Update(form *domain.Form) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE forms SET status = $1 WHERE id = $2`

	res, err := f.DB.ExecContext(ctx, query, form.Status, form.ID)
	if err != nil {
		return err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return domain.ErrFormNotFound
	}

	return nil
}
