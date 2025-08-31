package main

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	"log"
	"net/http"
	"strings"
)

type ctxKey string

const (
	userIDKey   ctxKey = "userID"
	userRoleKey ctxKey = "userRole"
)

func (app *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			WriteJSONError(w, http.StatusUnauthorized, errors.New("authorization header is required"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			WriteJSONError(w, http.StatusUnauthorized, errors.New("bearer token is required"))

			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(app.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			WriteJSONError(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			WriteJSONError(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		userIDStr, ok1 := claims["id"].(string)
		userRoleStr, ok2 := claims["role"].(string)
		if !ok1 || !ok2 {
			WriteJSONError(w, http.StatusUnauthorized, errors.New("invalid token payload"))
			return
		}

		userRole := domain.Role(userRoleStr)

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			WriteJSONError(w, http.StatusUnauthorized, errors.New("invalid user ID format"))
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, userRoleKey, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) RequireRoles(allowedRoles ...domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := GetUserRole(r.Context())
			if !ok {
				WriteJSONError(w, http.StatusUnauthorized, errors.New("missing role"))
				return
			}

			for _, role := range allowedRoles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			WriteJSONError(w, http.StatusForbidden, errors.New("permission denied"))
		})
	}
}

func (app *application) RequireActiveStatus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r.Context())
		if !ok {
			WriteJSONError(w, http.StatusUnauthorized, errors.New("missing user ID"))
			return
		}

		isActive, err := app.repositories.Users.IsActive(userID)
		if err != nil {
			WriteJSONError(w, http.StatusNotFound, err)
			log.Printf("failed to get isActive status: %V", err)
			return
		}

		if !isActive {
			WriteJSONError(w, http.StatusForbidden, errors.New("account is not active"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
