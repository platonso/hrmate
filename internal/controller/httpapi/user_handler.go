package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/controller/httpapi/dto"
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/service"
	"net/http"
	"time"
)

type UserHandler struct {
	userService *service.UserService
	validator   *validator.Validate
}

func NewUsersHandler(userService *service.UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   v,
	}
}

func (h *UserHandler) HandleUpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	var req dto.UserStatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	user, err := h.userService.UpdateStatus(ctx, userID, req.Status)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			dto.WriteJSONError(w, http.StatusNotFound, err)
		default:
			dto.WriteJSONError(w, http.StatusInternalServerError, nil)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	role, ok := GetUserRole(ctx)
	if !ok {
		dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("missing role"))
		return
	}

	users, err := h.userService.GetUsersByRole(ctx, role)
	if err != nil {
		dto.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}
