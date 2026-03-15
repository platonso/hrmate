package form

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	"github.com/platonso/hrmate/internal/handler/form/dto"
)

type Repository interface {
	Create(ctx context.Context, form *domain.Form) error
	FindByFormID(ctx context.Context, formId uuid.UUID) (*domain.Form, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Form, error)
	FindByUserIDWithUser(ctx context.Context, userID uuid.UUID) ([]dto.FormsWithUserResponse, error)
	FindAllWithUsers(ctx context.Context) ([]dto.FormsWithUserResponse, error)
	Update(ctx context.Context, form *domain.Form) error
}

type UserRepository interface {
	FindByUserID(ctx context.Context, userId uuid.UUID) (*domain.User, error)
}

type Service struct {
	formRepo Repository
	userRepo UserRepository
}

func NewService(formRepo Repository, userRepo UserRepository) *Service {
	return &Service{
		formRepo: formRepo,
		userRepo: userRepo,
	}
}

func (s *Service) Create(ctx context.Context, formDTO *dto.FormCreateRequest, userID uuid.UUID) (*domain.Form, error) {
	form := domain.NewForm(
		userID,
		formDTO.Title,
		formDTO.Description,
		formDTO.StartDate,
		formDTO.EndDate,
	)

	if err := s.formRepo.Create(ctx, &form); err != nil {
		return nil, fmt.Errorf("create form: %w", err)
	}

	return &form, nil
}

func (s *Service) GetForm(ctx context.Context, formID, currentUserID uuid.UUID) (*dto.FormWithUserResponse, error) {

	currentUser, err := s.userRepo.FindByUserID(ctx, currentUserID)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	form, err := s.formRepo.FindByFormID(ctx, formID)
	if err != nil {
		if errors.Is(err, errs.ErrFormNotFound) {
			return nil, errs.ErrFormNotFound
		}
		return nil, fmt.Errorf("find form: %w", err)
	}

	author, err := s.userRepo.FindByUserID(ctx, form.UserID)
	if err != nil {
		return nil, fmt.Errorf("find userId: %w", err)
	}

	formWithUser := dto.FormWithUserResponse{
		User: *author,
		Form: *form,
	}

	switch currentUser.Role {
	case domain.RoleAdmin, domain.RoleHR:
		return &formWithUser, nil

	case domain.RoleEmployee:
		if form.UserID == currentUser.ID {
			return &formWithUser, nil
		}
	}
	return nil, errs.ErrForbidden
}

func (s *Service) GetAllForms(ctx context.Context, currentUserID uuid.UUID) ([]dto.FormsWithUserResponse, error) {
	currentUser, err := s.userRepo.FindByUserID(ctx, currentUserID)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	switch currentUser.Role {
	case domain.RoleAdmin, domain.RoleHR:
		formsWithUsers, err := s.formRepo.FindAllWithUsers(ctx)
		if err != nil {
			return nil, err
		}
		return formsWithUsers, nil

	case domain.RoleEmployee:
		formsWithUsers, err := s.formRepo.FindByUserIDWithUser(ctx, currentUserID)
		if err != nil {
			return nil, err
		}
		if len(formsWithUsers) == 0 {
			return []dto.FormsWithUserResponse{}, nil
		}
		if formsWithUsers[0].User.ID == currentUser.ID {
			return formsWithUsers, nil
		}
	}

	return nil, errs.ErrForbidden
}

func (s *Service) Update(ctx context.Context, newStatus *dto.FormStatusUpdateRequest, formID uuid.UUID) (*domain.Form, error) {
	form, err := s.formRepo.FindByFormID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("find form: %w", err)
	}

	if newStatus.Status != form.Status {
		form.UpdateStatus(newStatus.Status)

		if err := s.formRepo.Update(ctx, form); err != nil {
			return nil, fmt.Errorf("update form: %w", err)
		}
	}

	return form, nil
}
