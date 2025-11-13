package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/controller/httpapi/dto"
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/repository"
)

type FormService struct {
	formRepo repository.FormRepository
	userRepo repository.UserRepository
}

func NewFormService(
	formRepo repository.FormRepository,
	userRepo repository.UserRepository,
) *FormService {
	return &FormService{
		formRepo: formRepo,
		userRepo: userRepo,
	}
}

func (s *FormService) Create(ctx context.Context, formDTO *dto.FormCreateRequest, userID uuid.UUID) (*domain.Form, error) {
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

func (s *FormService) GetForm(ctx context.Context, formID, currentUserID uuid.UUID) (*dto.FormWithUserResponse, error) {

	currentUser, err := s.userRepo.FindByUserId(ctx, currentUserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	form, err := s.formRepo.FindByFormID(ctx, formID)
	if err != nil {
		if errors.Is(err, domain.ErrFormNotFound) {
			return nil, domain.ErrFormNotFound
		}
		return nil, fmt.Errorf("find form: %w", err)
	}

	author, err := s.userRepo.FindByUserId(ctx, form.UserID)
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
	return nil, domain.ErrForbidden
}

func (s *FormService) GetAllForms(ctx context.Context, currentUserID uuid.UUID) ([]dto.FormsWithUserResponse, error) {
	currentUser, err := s.userRepo.FindByUserId(ctx, currentUserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
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
		if formsWithUsers[0].User.ID == currentUser.ID {
			return formsWithUsers, nil
		}
	}

	return nil, domain.ErrForbidden
}

func (s *FormService) Update(ctx context.Context, newStatus *dto.FormStatusUpdateRequest, formID, currentUserID uuid.UUID) (*domain.Form, error) {
	form, err := s.formRepo.FindByFormID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("find form: %w", err)
	}

	form.UpdateStatus(newStatus.Status)

	if err := s.formRepo.Update(ctx, form); err != nil {
		return nil, fmt.Errorf("update form: %w", err)
	}

	return form, nil
}
