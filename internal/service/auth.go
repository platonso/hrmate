package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/controller/httpapi/dto"
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

type AuthService struct {
	userRepo  repository.UserRepository
	JWTSecret string
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) ImplementAdmin(ctx context.Context) error {
	admin, err := s.userRepo.FindAdmin(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return fmt.Errorf("admin user not found")
		}
		return fmt.Errorf("failed to find admin")
	}

	if admin != nil {
		return nil
	}

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || password == "" {
		return errors.New("ADMIN_EMAIL and ADMIN_PASSWORD must be set in environment")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	adminUser := domain.User{
		ID:             uuid.New(),
		Role:           domain.RoleAdmin,
		FirstName:      "Super",
		LastName:       "User",
		Position:       "Administrator",
		Email:          email,
		HashedPassword: string(hashedPassword),
		IsActive:       true,
	}

	if err := s.userRepo.Create(ctx, &adminUser); err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	log.Println("Admin created")
	return nil

}

func (s *AuthService) Register(ctx context.Context, userDTO *dto.RegisterRequest) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	var user domain.User
	user = domain.NewUser(
		userDTO.Role,
		userDTO.FirstName,
		userDTO.LastName,
		userDTO.Position,
		userDTO.Email,
		string(hashedPassword),
	)

	if user.Role == domain.RoleEmployee {
		user.ChangeStatus(true)
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}

	token, err := generateJWT(user.ID, user.Role, s.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("generate jwt: %w", err)
	}

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, userDTO *dto.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, userDTO.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", fmt.Errorf("find user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userDTO.Password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := generateJWT(user.ID, user.Role, s.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("generate jwt: %w", err)
	}

	return token, nil
}

func generateJWT(userID uuid.UUID, role domain.Role, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   userID,
		"role": role,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
