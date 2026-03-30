package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/errors"
)

type Repository interface {
	Update(ctx context.Context, user *domain.User) error
	FindByRole(ctx context.Context, roles ...domain.Role) ([]domain.User, error)
	FindByUserID(ctx context.Context, userId uuid.UUID) (*domain.User, error)
	IsActive(ctx context.Context, userID uuid.UUID) (bool, error)
}
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

func (s *Service) ChangeActiveStatus(ctx context.Context, userID uuid.UUID, isActive bool) (*domain.User, error) {
	user, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	var changed bool
	if isActive {
		changed = user.Activate()
	} else {
		changed = user.Deactivate()
	}

	if changed {
		if err := s.repo.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("update user: %w", err)
		}
	}

	return user, nil
}

func (s *Service) GetUsersByRole(ctx context.Context, requesterRole domain.Role) ([]domain.User, error) {
	var rolesToQuery []domain.Role

	switch requesterRole {
	case domain.RoleAdmin:
		rolesToQuery = []domain.Role{domain.RoleHR, domain.RoleEmployee}
	case domain.RoleHR:
		rolesToQuery = []domain.Role{domain.RoleEmployee}
	default:
		return nil, errors.ErrForbidden
	}

	users, err := s.repo.FindByRole(ctx, rolesToQuery...)
	if err != nil {
		return nil, fmt.Errorf("find users by role: %w", err)
	}

	if users == nil {
		users = []domain.User{}
	}

	return users, nil
}

func (s *Service) IsActive(ctx context.Context, userID uuid.UUID) (bool, error) {
	active, err := s.repo.IsActive(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("check user active status: %w", err)
	}
	return active, nil
}
