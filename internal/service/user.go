package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

type UpdateUserStatusCommand struct {
	UserID uuid.UUID
	Status bool
}

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

func (s *UserService) UpdateStatus(ctx context.Context, userID uuid.UUID, newStatus bool) (*domain.User, error) {
	user, err := s.userRepo.FindByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	user.ChangeStatus(newStatus)

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUsersByRole(ctx context.Context, requesterRole domain.Role) ([]domain.User, error) {
	var rolesToQuery []domain.Role

	switch requesterRole {
	case domain.RoleAdmin:
		rolesToQuery = []domain.Role{domain.RoleHR, domain.RoleEmployee}
	case domain.RoleHR:
		rolesToQuery = []domain.Role{domain.RoleEmployee}
	default:
		return nil, domain.ErrForbidden
	}

	users, err := s.userRepo.FindAllByRole(ctx, rolesToQuery...)
	if err != nil {
		return nil, fmt.Errorf("find users by role: %w", err)
	}

	if users == nil {
		users = []domain.User{}
	}

	return users, nil
}

func (s *UserService) IsActive(ctx context.Context, userID uuid.UUID) (bool, error) {
	active, err := s.userRepo.IsActive(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("check user active status: %w", err)
	}
	return active, nil
}
