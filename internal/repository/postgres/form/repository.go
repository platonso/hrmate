package form

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	"github.com/platonso/hrmate/internal/handler/form/dto"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (f *Repository) Create(ctx context.Context, form *domain.Form) error {
	query := `
		INSERT INTO forms (id, user_id, title, description, start_date, end_date, created_at, approved_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`
	_, err := f.DB.Exec(
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

func (f *Repository) FindByFormID(ctx context.Context, formId uuid.UUID) (*domain.Form, error) {
	query := `
		SELECT 
		    f.id, f.user_id, f.title, f.description, f.start_date, 
		    f.end_date, f.created_at, f.approved_at, f.status
		FROM forms f 
		JOIN users u ON u.id = f.user_id
		WHERE f.id = $1
`

	var form domain.Form

	err := f.DB.QueryRow(ctx, query, formId).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrFormNotFound
		}
		return nil, err
	}

	return &form, nil
}

func (f *Repository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Form, error) {
	formsQuery := `
        SELECT id, user_id, title, description, start_date, end_date, created_at, approved_at, status
        FROM forms WHERE user_id = $1
    `
	rows, err := f.DB.Query(ctx, formsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []domain.Form
	for rows.Next() {
		var form domain.Form
		err := rows.Scan(
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

	return forms, nil
}

func (f *Repository) FindByUserIDWithUser(ctx context.Context, userID uuid.UUID) ([]dto.FormsWithUserResponse, error) {

	// TODO: fix

	query := `
		SELECT 
		    u.id, u.first_name, u.last_name, u.position, u.email, u.is_active,
		    f.id, f.user_id, f.title, f.description, f.start_date, 
		    f.end_date, f.created_at, f.approved_at, f.status
		FROM forms f 
		JOIN users u ON u.id = f.user_id
		WHERE u.id = $1
		ORDER BY f.created_at DESC
`

	rows, err := f.DB.Query(ctx, query, userID)
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
			&user.IsActive,
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
		return []dto.FormsWithUserResponse{}, nil
	}

	result := []dto.FormsWithUserResponse{
		{
			User:  user,
			Forms: forms,
		},
	}

	return result, nil
}

func (f *Repository) FindAllWithUsers(ctx context.Context) ([]dto.FormsWithUserResponse, error) {

	query := `
		SELECT 
			u.id, u.first_name, u.last_name, u.position, u.email, u.is_active,
		    f.id, f.user_id, f.title, f.description, f.start_date, 
		    f.end_date, f.created_at, f.approved_at, f.status
		FROM users u
		JOIN forms f ON u.id = f.user_id
		ORDER BY f.created_at DESC
`
	rows, err := f.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersMap = make(map[uuid.UUID]*dto.FormsWithUserResponse)

	for rows.Next() {
		var user domain.User
		var form domain.Form

		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Position, &user.Email, &user.IsActive,
			&form.ID, &form.UserID, &form.Title, &form.Description, &form.StartDate, &form.EndDate,
			&form.CreatedAt, &form.ApprovedAt, &form.Status,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := usersMap[user.ID]; !ok {
			usersMap[user.ID] = &dto.FormsWithUserResponse{
				User:  user,
				Forms: []domain.Form{},
			}
		}

		usersMap[user.ID].Forms = append(usersMap[user.ID].Forms, form)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var result []dto.FormsWithUserResponse
	for _, u := range usersMap {
		result = append(result, *u)
	}

	return result, nil
}

func (f *Repository) Update(ctx context.Context, form *domain.Form) error {
	query := `
	UPDATE forms 
	SET 
	    status = $1,
	    approved_at = $2
		
	WHERE id = $3`

	tag, err := f.DB.Exec(ctx, query, form.Status, form.ApprovedAt, form.ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrFormNotFound
	}

	return nil
}
