package form

import (
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
)

type Filter struct {
	UserID     *uuid.UUID
	FormStatus *domain.FormStatus
}

func (f *Filter) ValidateStatus() error {
	if f.FormStatus == nil {
		return nil
	}
	switch *f.FormStatus {
	case domain.StatusPending, domain.StatusApproved, domain.StatusRejected:
		return nil
	default:
		return errs.ErrFormInvalidStatus
	}
}
