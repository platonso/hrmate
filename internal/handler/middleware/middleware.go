package middleware

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
	"github.com/platonso/hrmate/internal/handler/middleware/dto"
)

const (
	userIDKey   = "userID"
	userRoleKey = "userRole"
)

type AuthService interface {
	GetJWTSecret() string
}

type UserService interface {
	IsActive(ctx context.Context, userID uuid.UUID) (bool, error)
}

type Auth struct {
	AuthSvc AuthService
	UserSvc UserService
}

func (m *Auth) AuthMiddleware(next http.Handler) http.Handler {
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

func (m *Auth) RequireRoles(allowedRoles ...domain.Role) func(http.Handler) http.Handler {
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

func (m *Auth) RequireActiveStatus(next http.Handler) http.Handler {
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
