package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	authdto "github.com/platonso/hrmate/internal/handler/auth/dto"
	"github.com/platonso/hrmate/internal/handler/middleware"
	"github.com/platonso/hrmate/internal/handler/middleware/dto"
)

type Service interface {
	GetUsersByRole(ctx context.Context, requesterRole domain.Role) ([]domain.User, error)
	UpdateStatus(ctx context.Context, userID uuid.UUID, newStatus bool) (*domain.User, error)
	//GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.Username, error)
	//IsActive(ctx context.Context, userID uuid.UUID) (bool, error)
}

type Handler struct {
	svc       Service
	validator *validator.Validate
}

func NewHandler(svc Service, v *validator.Validate) *Handler {
	return &Handler{
		svc:       svc,
		validator: v,
	}
}

func (h *Handler) HandleUpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	var req authdto.UserStatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	user, err := h.svc.UpdateStatus(r.Context(), userID, req.Status)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			dto.WriteJSONError(w, http.StatusNotFound, err)
		default:
			dto.WriteJSONError(w, http.StatusInternalServerError, errors.New("internal server error"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func (h *Handler) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	role, ok := middleware.GetUserRole(ctx)
	if !ok {
		dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("missing role"))
		return
	}

	users, err := h.svc.GetUsersByRole(ctx, role)
	if err != nil {
		dto.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}
