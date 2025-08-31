package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

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

func (app *application) HandleRegister(role domain.Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerRequestDTO RegisterRequestDTO

		if err := json.NewDecoder(r.Body).Decode(&registerRequestDTO); err != nil {
			WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
			return
		}

		if err := app.validator.Struct(registerRequestDTO); err != nil {
			WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequestDTO.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("failed to hash password: %v", err)
			return
		}

		var user domain.User
		switch role {
		case domain.RoleEmployee:
			user = domain.NewEmployee(
				registerRequestDTO.FirstName,
				registerRequestDTO.LastName,
				registerRequestDTO.Position,
				registerRequestDTO.Email,
				string(hashedPassword),
			)
		case domain.RoleHR:
			user = domain.NewHR(
				registerRequestDTO.FirstName,
				registerRequestDTO.LastName,
				registerRequestDTO.Position,
				registerRequestDTO.Email,
				string(hashedPassword),
			)
		default:
			WriteJSONError(w, http.StatusBadRequest, errors.New("invalid role"))
			return
		}

		if err := app.repositories.Users.Create(&user); err != nil {
			WriteJSONError(w, http.StatusConflict, domain.ErrUserAlreadyExists)
			return
		}

		tokenString, err := generateJWT(user.ID, user.Role, app.jwtSecret)
		if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, err)
			log.Printf("error generating token: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		registerResponseDTO := AuthResponseDTO{Token: tokenString}
		if err := json.NewEncoder(w).Encode(registerResponseDTO); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	}

}

func (app *application) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequestDTO LoginRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&loginRequestDTO); err != nil {
		WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := app.validator.Struct(loginRequestDTO); err != nil {
		WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	user, err := app.repositories.Users.FindByEmail(loginRequestDTO.Email)
	if err != nil {
		log.Printf("Database error finding user by email: %v", err)
		WriteJSONError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}

	if user == nil {
		WriteJSONError(w, http.StatusUnauthorized, errors.New("invalid email or password"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginRequestDTO.Password))
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, errors.New("invalid email or password"))
		return
	}

	tokenString, err := generateJWT(user.ID, user.Role, app.jwtSecret)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, errors.New("error generating token"))
		log.Printf("error generating token: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	loginResponseDTO := AuthResponseDTO{Token: tokenString}
	if err := json.NewEncoder(w).Encode(loginResponseDTO); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}
