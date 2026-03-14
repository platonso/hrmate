package httpapi

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	authdto "github.com/platonso/hrmate/internal/handler/auth/dto"
	"github.com/platonso/hrmate/internal/handler/httpapi/dto"
)

const (
	userIDKey   = "userID"
	userRoleKey = "userRole"
)

type AuthService interface {
	Register(ctx context.Context, userDTO *authdto.RegisterRequest) (string, error)
	Login(ctx context.Context, userDTO *authdto.LoginRequest) (string, error)
	GetJWTSecret() string
}

type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateStatus(ctx context.Context, userID uuid.UUID, newStatus bool) (*domain.User, error)
	GetUsersByRole(ctx context.Context, requesterRole domain.Role) ([]domain.User, error)
	IsActive(ctx context.Context, userID uuid.UUID) (bool, error)
}

type AuthMiddleware struct {
	AuthSvc AuthService
	UserSvc UserService
}

func (m *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("authorization header is required"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("bearer token is required"))

			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.AuthSvc.GetJWTSecret()), nil
		})

		if err != nil || !token.Valid {
			dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		userIDStr, ok1 := claims["id"].(string)
		userRoleStr, ok2 := claims["role"].(string)
		if !ok1 || !ok2 {
			dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("invalid token payload"))
			return
		}

		userRole := domain.Role(userRoleStr)

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("invalid user ID format"))
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, userRoleKey, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRoles(allowedRoles ...domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := GetUserRole(r.Context())
			if !ok {
				dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("missing role"))
				return
			}

			for _, role := range allowedRoles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			dto.WriteJSONError(w, http.StatusForbidden, errors.New("permission denied"))
		})
	}
}

//func (m *AuthMiddleware) RequireRoles(allowedRoles ...domain.Role) func(httpapi.Handler) httpapi.Handler {
//	return func(next httpapi.Handler) httpapi.Handler {
//		return httpapi.HandlerFunc(func(w httpapi.ResponseWriter, r *httpapi.Request) {
//			userID, ok := GetUserID(r.Context())
//			if !ok {
//				dto.WriteJSONError(w, httpapi.StatusUnauthorized, errors.New("missing user ID"))
//				return
//			}
//
//			user, err := m.userService.GetUserByID(r.Context(), userID)
//			if err != nil {
//				dto.WriteJSONError(w, httpapi.StatusUnauthorized, domain.ErrInvalidCredentials)
//				log.Printf("failed to get user in middleware: %v", err)
//				return
//			}
//
//			for _, role := range allowedRoles {
//				if user.Role == role {
//					next.ServeHTTP(w, r)
//					return
//				}
//			}
//
//			dto.WriteJSONError(w, httpapi.StatusForbidden, errors.New("permission denied"))
//		})
//	}
//}

func (m *AuthMiddleware) RequireActiveStatus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r.Context())
		if !ok {
			dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("missing user ID"))
			return
		}

		isActive, err := m.UserSvc.IsActive(r.Context(), userID)
		if err != nil {
			dto.WriteJSONError(w, http.StatusUnauthorized, errs.ErrInvalidCredentials)
			log.Printf("failed to get isActive status: %v", err)
			return
		}

		if !isActive {
			dto.WriteJSONError(w, http.StatusForbidden, errors.New("account is not active"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
