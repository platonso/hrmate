package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	errs "github.com/platonso/hrmate/internal/errors"
	authdto "github.com/platonso/hrmate/internal/handler/auth/dto"
	errdto "github.com/platonso/hrmate/internal/handler/middleware/dto"
	"github.com/platonso/hrmate/internal/service/auth/model"
)

type Service interface {
	Register(ctx context.Context, registerInput *model.RegisterInput) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
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

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req authdto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errdto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		errdto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	token, err := h.svc.Register(r.Context(), authdto.ToRegisterInput(&req))
	if err != nil {
		errdto.WriteJSONError(w, http.StatusConflict, errs.ErrUserAlreadyExists)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(authdto.AuthResponse{Token: token})

}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req authdto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errdto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		errdto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	token, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidCredentials):
			errdto.WriteJSONError(w, http.StatusUnauthorized, err)
		case errors.Is(err, errs.ErrUserNotActive):
			errdto.WriteJSONError(w, http.StatusForbidden, err)
		default:
			errdto.WriteJSONError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authdto.AuthResponse{Token: token})
}
