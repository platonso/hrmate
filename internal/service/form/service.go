package form

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	"github.com/platonso/hrmate/internal/service/form/model"
)

type Repository interface {
	Create(ctx context.Context, form *domain.Form) error
	FindAll(ctx context.Context) ([]domain.Form, error)
	FindByFormID(ctx context.Context, formId uuid.UUID) (*domain.Form, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Form, error)
	Update(ctx context.Context, form *domain.Form) error
}

type UserRepository interface {
	FindByUserID(ctx context.Context, userId uuid.UUID) (*domain.User, error)
	FindByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]domain.User, error)
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

func (s *Service) Create(ctx context.Context, formInput *model.FormCreateInput, userID uuid.UUID) (*domain.Form, error) {
	form := domain.NewForm(
		userID,
		formInput.Title,
		formInput.Description,
		formInput.StartDate,
		formInput.EndDate,
	)

	if err := s.formRepo.Create(ctx, &form); err != nil {
		return nil, fmt.Errorf("create form: %w", err)
	}

	return &form, nil
}

func (s *Service) GetForm(ctx context.Context, formID, currentUserID uuid.UUID) (*domain.Form, error) {

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

	switch currentUser.Role {
	case domain.RoleAdmin, domain.RoleHR:
		return form, nil

	case domain.RoleEmployee:
		if form.UserID == currentUser.ID {
			return form, nil
		}
	}
	return nil, errs.ErrForbidden
}

func (s *Service) GetForms(ctx context.Context, targetUserID, requesterUserID uuid.UUID) ([]domain.Form, error) {
	//requester, err := s.userRepo.FindByUserID(ctx, requesterUserID)
	//if err != nil {
	//	return nil, errs.ErrUserNotFound
	//}
	
	targetUser, err := s.userRepo.FindByUserID(ctx, targetUserID)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	forms, err := s.formRepo.FindByUserID(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find forms for user %s: %w", targetUserID, err)
	}

	if len(forms) == 0 {
		return []model.FormsWithUser{}, nil
	}

	result := []model.FormsWithUser{
		{
			User:  *targetUser,
			Forms: forms,
		},
	}

	return result, nil
}

//func (s *Service) GetFormsWithUser(ctx context.Context, targetUserID, requesterUserID uuid.UUID) ([]model.FormsWithUser, error) {
//	requester, err := s.userRepo.FindByUserID(ctx, requesterUserID)
//	if err != nil {
//		return nil, errs.ErrUserNotFound
//	}
//
//	if requester.Role != domain.RoleHR && requester.Role != domain.RoleAdmin {
//		return nil, errs.ErrForbidden
//	}
//
//	targetUser, err := s.userRepo.FindByUserID(ctx, targetUserID)
//	if err != nil {
//		return nil, errs.ErrUserNotFound
//	}
//
//	forms, err := s.formRepo.FindByUserID(ctx, targetUserID)
//	if err != nil {
//		return nil, fmt.Errorf("failed to find forms for user %s: %w", targetUserID, err)
//	}
//
//	if len(forms) == 0 {
//		return []model.FormsWithUser{}, nil
//	}
//
//	result := []model.FormsWithUser{
//		{
//			User:  *targetUser,
//			Forms: forms,
//		},
//	}
//
//	return result, nil
//}

func (s *Service) GetFormsWithUsers(ctx context.Context, requesterUserID uuid.UUID) ([]model.FormsWithUser, error) {
	// Get current user
	currentUser, err := s.userRepo.FindByUserID(ctx, requesterUserID)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	var forms []domain.Form
	var result []model.FormsWithUser

	switch currentUser.Role {
	case domain.RoleAdmin, domain.RoleHR:
		// For admin and HR - all forms
		forms, err = s.formRepo.FindAll(ctx)
		if err != nil {
			return nil, err
		}

		// Collect unique form's authors
		authorIDs := make([]uuid.UUID, 0, len(forms))
		seen := make(map[uuid.UUID]struct{})

		for _, form := range forms {
			if _, ok := seen[form.UserID]; !ok {
				seen[form.UserID] = struct{}{}
				authorIDs = append(authorIDs, form.UserID)
			}
		}

		// Get authors
		authors, err := s.userRepo.FindByUserIDs(ctx, authorIDs)
		if err != nil {
			return nil, err
		}

		// Map for authors
		authorMap := make(map[uuid.UUID]domain.User)
		for _, author := range authors {
			authorMap[author.ID] = author
		}

		result = make([]model.FormsWithUser, 0, len(forms))
		for _, form := range forms {
			if author, ok := authorMap[form.UserID]; ok {
				result = append(result, model.FormsWithUser{
					User:  author,
					Forms: []domain.Form{form},
				})
			}
		}

	case domain.RoleEmployee:
		// for employee - only his forms
		forms, err = s.formRepo.FindByUserID(ctx, requesterUserID)
		if err != nil {
			return nil, err
		}

		if len(forms) == 0 {
			return []model.FormsWithUser{}, nil
		}

		result = []model.FormsWithUser{
			{
				User:  *currentUser,
				Forms: forms,
			},
		}
	}

	return result, nil
}

func (s *Service) Approve(ctx context.Context, formID uuid.UUID, comment string) (*domain.Form, error) {
	form, err := s.formRepo.FindByFormID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("find form: %w", err)
	}

	if form.Status != domain.StatusApproved {
		form.ApproveForm(comment)

		if err := s.formRepo.Update(ctx, form); err != nil {
			return nil, fmt.Errorf("update form: %w", err)
		}
	}

	return form, nil
}
