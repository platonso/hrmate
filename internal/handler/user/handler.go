package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	"github.com/platonso/hrmate/internal/handler/middleware"
	"github.com/platonso/hrmate/internal/handler/middleware/dto"
	userdto "github.com/platonso/hrmate/internal/handler/user/dto"
)

type Service interface {
	GetUsersByRole(ctx context.Context, requesterRole domain.Role) ([]domain.User, error)
	ChangeActiveStatus(ctx context.Context, userID uuid.UUID, newStatus bool) (*domain.User, error)
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

func (h *Handler) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	requesterRole, ok := middleware.GetUserRole(r.Context())
	if !ok {
		dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("missing role in context"))
		return
	}

	users, err := h.svc.GetUsersByRole(r.Context(), requesterRole)
	if err != nil {
		dto.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(userdto.ToUserResponses(users))
}

func (h *Handler) HandleActivate(w http.ResponseWriter, r *http.Request) {
	h.handleChangeActiveStatus(w, r, true)
}
func (h *Handler) HandleDeactivate(w http.ResponseWriter, r *http.Request) {
	h.handleChangeActiveStatus(w, r, false)
}

func (h *Handler) handleChangeActiveStatus(
	w http.ResponseWriter,
	r *http.Request,
	newStatus bool,
) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}
	user, err := h.svc.ChangeActiveStatus(r.Context(), userID, newStatus)
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
	_ = json.NewEncoder(w).Encode(userdto.ToUserResponse(user))
}
