package entity

import "github.com/platonso/hrmate/internal/domain"

func ToFormRecord(f domain.Form) FormRecord {
	return FormRecord{
		ID:          f.ID,
		UserID:      f.UserID,
		Title:       f.Title,
		Description: f.Description,

		StartDate:  f.StartDate,
		EndDate:    f.EndDate,
		CreatedAt:  f.CreatedAt,
		ReviewedAt: f.ReviewedAt,

		Status:  string(f.Status),
		Comment: f.Comment,
	}
}

func ToDomainForm(fr FormRecord) domain.Form {
	return domain.Form{
		ID:          fr.ID,
		UserID:      fr.UserID,
		Title:       fr.Title,
		Description: fr.Description,

		StartDate:  fr.StartDate,
		EndDate:    fr.EndDate,
		CreatedAt:  fr.CreatedAt,
		ReviewedAt: fr.ReviewedAt,

		Status:  domain.FormStatus(fr.Status),
		Comment: fr.Comment,
	}
}

func ToDomainForms(records []FormRecord) []domain.Form {
	forms := make([]domain.Form, len(records))
	for i := range records {
		forms[i] = ToDomainForm(records[i])
	}
	return forms
}
