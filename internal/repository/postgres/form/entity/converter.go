package entity

import (
	"database/sql"

	"github.com/platonso/hrmate/internal/domain"
)

func ToFormRecord(f domain.Form) FormRecord {
	record := FormRecord{
		ID:          f.ID,
		UserID:      f.UserID,
		Title:       f.Title,
		Description: f.Description,
		CreatedAt:   f.CreatedAt,
		Status:      string(f.Status),
	}

	if f.StartDate != nil {
		record.StartDate = sql.NullTime{
			Time:  *f.StartDate,
			Valid: true,
		}
	} else {
		record.StartDate = sql.NullTime{Valid: false}
	}

	if f.EndDate != nil {
		record.EndDate = sql.NullTime{
			Time:  *f.EndDate,
			Valid: true,
		}
	} else {
		record.EndDate = sql.NullTime{Valid: false}
	}

	if f.ApprovedAt != nil {
		record.ApprovedAt = sql.NullTime{
			Time:  *f.ApprovedAt,
			Valid: true,
		}
	} else {
		record.ApprovedAt = sql.NullTime{Valid: false}
	}

	if f.Comment != nil {
		record.Comment = sql.NullString{
			String: *f.Comment,
			Valid:  true,
		}
	} else {
		record.Comment = sql.NullString{Valid: false}
	}

	return record
}

func ToDomainForm(record FormRecord) domain.Form {
	form := domain.Form{
		ID:          record.ID,
		UserID:      record.UserID,
		Title:       record.Title,
		Description: record.Description,
		CreatedAt:   record.CreatedAt,
		Status:      domain.FormStatus(record.Status),
	}

	if record.StartDate.Valid {
		form.StartDate = &record.StartDate.Time
	}

	if record.EndDate.Valid {
		form.EndDate = &record.EndDate.Time
	}

	if record.ApprovedAt.Valid {
		form.ApprovedAt = &record.ApprovedAt.Time
	}

	if record.Comment.Valid {
		form.Comment = &record.Comment.String
	}

	return form
}

func ToDomainForms(records []FormRecord) []domain.Form {
	if len(records) == 0 {
		return []domain.Form{}
	}

	forms := make([]domain.Form, 0, len(records))
	for _, r := range records {
		form := ToDomainForm(r)
		forms = append(forms, form)
	}
	return forms
}
