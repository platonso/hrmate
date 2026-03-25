package form

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	"github.com/platonso/hrmate/internal/repository/postgres/form/entity"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(ctx context.Context, form *domain.Form) error {
	rec := entity.ToFormRecord(*form)
	query := `
		INSERT INTO forms (id, user_id, title, description, start_date, end_date, created_at, approved_at, status, comment)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`
	_, err := r.DB.Exec(
		ctx,
		query,
		rec.ID,
		rec.UserID,
		rec.Title,
		rec.Description,
		rec.StartDate,
		rec.EndDate,
		rec.CreatedAt,
		rec.ApprovedAt,
		rec.Status,
		rec.Comment,
	)
	return err
}

func (r *Repository) FindAll(ctx context.Context) ([]domain.Form, error) {
	query := `
        SELECT id, user_id, title, description, start_date, end_date, created_at, approved_at, status, comment
        FROM forms
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query forms: %w", err)
	}
	defer rows.Close()

	records, err := scanForms(rows)
	if err != nil {
		return nil, err
	}

	return entity.ToDomainForms(records), nil
}

func (r *Repository) FindByFormID(ctx context.Context, formId uuid.UUID) (*domain.Form, error) {
	query := `
		SELECT id, user_id, title, description, start_date, end_date, created_at, approved_at, status, comment
		FROM forms
		WHERE id = $1
`
	rec, err := r.findForm(ctx, query, formId)
	if err != nil {
		return nil, err
	}
	form := entity.ToDomainForm(rec)
	return &form, nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Form, error) {
	query := `
        SELECT id, user_id, title, description, start_date, end_date, created_at, approved_at, status, comment
        FROM forms 
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query forms by user: %w", err)
	}
	defer rows.Close()

	records, err := scanForms(rows)
	if err != nil {
		return nil, err
	}

	return entity.ToDomainForms(records), nil
}

func (r *Repository) Update(ctx context.Context, form *domain.Form) error {
	rec := entity.ToFormRecord(*form)
	query := `
	UPDATE forms 
	SET approved_at = $1, status = $2, comment = $3
	WHERE id = $4`

	tag, err := r.DB.Exec(ctx, query, rec.ApprovedAt, rec.Status, rec.Comment, rec.ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrFormNotFound
	}

	return nil
}

func (r *Repository) findForm(ctx context.Context, query string, args ...any) (entity.FormRecord, error) {
	var rec entity.FormRecord
	err := r.DB.QueryRow(ctx, query, args...).Scan(
		&rec.ID,
		&rec.UserID,
		&rec.Title,
		&rec.Description,
		&rec.StartDate,
		&rec.EndDate,
		&rec.CreatedAt,
		&rec.ApprovedAt,
		&rec.Status,
		&rec.Comment,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.FormRecord{}, errs.ErrFormNotFound
		}
		return entity.FormRecord{}, err
	}
	return rec, nil
}

func scanForms(rows pgx.Rows) ([]entity.FormRecord, error) {
	var records []entity.FormRecord
	for rows.Next() {
		var rec entity.FormRecord
		err := rows.Scan(
			&rec.ID,
			&rec.UserID,
			&rec.Title,
			&rec.Description,
			&rec.StartDate,
			&rec.EndDate,
			&rec.CreatedAt,
			&rec.ApprovedAt,
			&rec.Status,
			&rec.Comment,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan form row: %w", err)
		}
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return records, nil
}
