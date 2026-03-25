package auth

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	"github.com/platonso/hrmate/internal/service/auth/model"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByRole(ctx context.Context, roles ...domain.Role) ([]domain.User, error)
}
type Service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ImplementAdmin(ctx context.Context, email, password string) error {

	admin, err := s.repo.FindByRole(ctx, domain.RoleAdmin)
	if err != nil {
		return fmt.Errorf("failed to find admin")
	}

	if len(admin) > 0 {
		log.Println("the existing admin is used")
		return nil
	}

	if email == "" || password == "" {
		return errors.New("ADMIN_EMAIL and ADMIN_PASSWORD must be set in environment")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	adminUser := domain.NewUser(
		domain.RoleAdmin,
		"Super",
		"UserRepository",
		"Administrator",
		email,
		string(hashedPassword),
	)

	adminUser.ChangeStatus(true)

	if err := s.repo.Create(ctx, &adminUser); err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	log.Println("admin has been created successfully")
	return nil
}

func (s *Service) Register(ctx context.Context, registerInput *model.RegisterInput) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	user := domain.NewUser(
		registerInput.Role,
		registerInput.FirstName,
		registerInput.LastName,
		registerInput.Position,
		registerInput.Email,
		string(hashedPassword),
	)

	if err := s.repo.Create(ctx, &user); err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}

	token, err := generateJWT(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("generate jwt: %w", err)
	}

	return token, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return "", errs.ErrInvalidCredentials
		}
		return "", fmt.Errorf("find user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", errs.ErrInvalidCredentials
	}

	token, err := generateJWT(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("generate jwt: %w", err)
	}

	return token, nil
}

func (s *Service) GetJWTSecret() string {
	return s.jwtSecret
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
