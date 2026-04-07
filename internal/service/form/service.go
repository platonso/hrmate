package form

import (
	"context"
	"errors"
	"log"

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
	FindByFilter(ctx context.Context, filter *Filter) ([]domain.Form, error)
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
		log.Printf("failed to create form for user %s: %v", userID, err)
		return nil, errs.ErrInternalServer
	}

	return &form, nil
}

func (s *Service) GetForm(ctx context.Context, formID uuid.UUID, requesterID uuid.UUID, requesterRole domain.Role) (*domain.Form, error) {
	form, err := s.formRepo.FindByFormID(ctx, formID)
	if err != nil {
		if errors.Is(err, errs.ErrFormNotFound) {
			return nil, errs.ErrFormNotFound
		}
		log.Printf("failed to find form: %v", err)
		return nil, errs.ErrInternalServer
	}

	switch requesterRole {
	case domain.RoleAdmin, domain.RoleHR:
		return form, nil

	case domain.RoleEmployee:
		if form.UserID == requesterID {
			return form, nil
		}
		return nil, errs.ErrFormNotFound
	}
	return nil, errs.ErrForbidden
}

// GetForms retrieves forms based on filter with access control
// For Employee: can only access own forms
// For HR/Admin: can access all forms with optional filtering
func (s *Service) GetForms(ctx context.Context, filter *Filter, requesterID uuid.UUID, requesterRole domain.Role) ([]domain.Form, error) {
	// Access control logic
	switch requesterRole {
	case domain.RoleEmployee:
		// Employee can only access own forms
		if filter.UserID != nil && *filter.UserID != requesterID {
			return nil, errs.ErrForbidden
		}
		// Force filter to requester's ID
		filter.UserID = &requesterID

	case domain.RoleHR, domain.RoleAdmin:
		// HR/Admin can access any forms
		// If user_id is specified, validate user exists
		if filter.UserID != nil {
			_, err := s.userRepo.FindByUserID(ctx, *filter.UserID)
			if err != nil {
				return nil, errs.ErrUserNotFound
			}
		}
	default:
		return nil, errs.ErrForbidden
	}

	// Validate status if provided
	if err := filter.ValidateStatus(); err != nil {
		return nil, err
	}

	// Fetch forms from repository
	forms, err := s.formRepo.FindByFilter(ctx, filter)
	if err != nil {
		log.Printf("failed to find forms: %v", err)
		return nil, errs.ErrInternalServer
	}

	// Return empty array if no results (not an error)
	if len(forms) == 0 {
		return []domain.Form{}, nil
	}

	return forms, nil
}

func (s *Service) GetFormsWithUsers(ctx context.Context, filter *Filter, requesterRole domain.Role) ([]model.FormsWithUser, error) {
	if requesterRole == domain.RoleEmployee {
		return nil, errs.ErrForbidden
	}

	// If user_id is specified, validate user exists
	if filter.UserID != nil {
		_, err := s.userRepo.FindByUserID(ctx, *filter.UserID)
		if err != nil {
			return nil, errs.ErrUserNotFound
		}
	}

	// Validate status if provided
	if err := filter.ValidateStatus(); err != nil {
		return nil, err
	}

	// Fetch forms with filter
	forms, err := s.formRepo.FindByFilter(ctx, filter)
	if err != nil {
		log.Printf("failed to find forms: %v", err)
		return nil, errs.ErrInternalServer
	}

	if len(forms) == 0 {
		return []model.FormsWithUser{}, nil
	}

	// Collect unique user IDs
	userIDsMap := make(map[uuid.UUID]struct{})
	for _, form := range forms {
		userIDsMap[form.UserID] = struct{}{}
	}

	userIDs := make([]uuid.UUID, 0, len(userIDsMap))
	for userID := range userIDsMap {
		userIDs = append(userIDs, userID)
	}

	// Fetch users
	users, err := s.userRepo.FindByUserIDs(ctx, userIDs)
	if err != nil {
		log.Printf("failed to find users: %v", err)
		return nil, errs.ErrInternalServer
	}

	// Create user map
	userMap := make(map[uuid.UUID]domain.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Group forms by user
	userFormsMap := make(map[uuid.UUID][]domain.Form)
	for _, form := range forms {
		userFormsMap[form.UserID] = append(userFormsMap[form.UserID], form)
	}

	// Build result
	result := make([]model.FormsWithUser, 0, len(userFormsMap))
	for userID, userForms := range userFormsMap {
		if user, ok := userMap[userID]; ok {
			result = append(result, model.FormsWithUser{
				User:  user,
				Forms: userForms,
			})
		}
	}

	return result, nil
}

func (s *Service) Approve(ctx context.Context, formID uuid.UUID, comment string) (*domain.Form, error) {
	form, err := s.formRepo.FindByFormID(ctx, formID)
	if err != nil {
		if errors.Is(err, errs.ErrFormNotFound) {
			return nil, errs.ErrFormNotFound
		}
		log.Printf("Failed to find form: %v", err)
		return nil, errs.ErrInternalServer
	}

	changed, err := form.ApproveForm(comment)
	if err != nil {
		return nil, err
	}

	if changed {
		if err := s.formRepo.Update(ctx, form); err != nil {
			log.Printf("Failed to update form: %v", err)
			return nil, errs.ErrInternalServer
		}
	}

	return form, nil
}

func (s *Service) Reject(ctx context.Context, formID uuid.UUID, comment string) (*domain.Form, error) {
	form, err := s.formRepo.FindByFormID(ctx, formID)
	if err != nil {
		if errors.Is(err, errs.ErrFormNotFound) {
			return nil, errs.ErrFormNotFound
		}
		log.Printf("Failed to find form: %v", err)
		return nil, errs.ErrInternalServer
	}

	changed, err := form.ApproveForm(comment)
	if err != nil {
		return nil, err
	}

	if changed {
		if err := s.formRepo.Update(ctx, form); err != nil {
			log.Printf("Failed to update form: %v", err)
			return nil, errs.ErrInternalServer
		}
	}

	return form, nil
}
