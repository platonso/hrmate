package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/handler/form/dto"
)

type User interface {
	Create(ctx context.Context, user *domain.User) error
	FindByUserId(ctx context.Context, userId uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindAdmin(ctx context.Context) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	FindAllByRole(ctx context.Context, roles ...domain.Role) ([]domain.User, error)
	IsActive(ctx context.Context, userID uuid.UUID) (bool, error)
}

type Form interface {
	Create(ctx context.Context, form *domain.Form) error
	FindByFormID(ctx context.Context, formId uuid.UUID) (*domain.Form, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Form, error)
	FindByUserIDWithUser(ctx context.Context, userID uuid.UUID) ([]dto.FormsWithUserResponse, error)
	FindAllWithUsers(ctx context.Context) ([]dto.FormsWithUserResponse, error)
	Update(ctx context.Context, form *domain.Form) error
}
