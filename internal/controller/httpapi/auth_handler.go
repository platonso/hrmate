package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/platonso/hrmate/internal/controller/httpapi/dto"
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/service"
	"net/http"
	"time"
)

type AuthHandler struct {
	authService *service.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(serviceAuth *service.AuthService, v *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: serviceAuth,
		validator:   v,
	}
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	token, err := h.authService.Register(ctx, &req)
	if err != nil {
		dto.WriteJSONError(w, http.StatusConflict, domain.ErrUserAlreadyExists)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(dto.AuthResponse{Token: token})

}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	token, err := h.authService.Login(ctx, &req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			dto.WriteJSONError(w, http.StatusUnauthorized, err)
		} else {
			dto.WriteJSONError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.AuthResponse{Token: token})
}
