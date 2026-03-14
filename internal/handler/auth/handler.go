package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	errs "github.com/platonso/hrmate/internal/errors"
	authdto "github.com/platonso/hrmate/internal/handler/auth/dto"
	errdto "github.com/platonso/hrmate/internal/handler/httpapi/dto"
)

type Service interface {
	Register(ctx context.Context, userDTO *authdto.RegisterRequest) (string, error)
	Login(ctx context.Context, userDTO *authdto.LoginRequest) (string, error)
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

	token, err := h.svc.Register(r.Context(), &req)
	if err != nil {
		errdto.WriteJSONError(w, http.StatusConflict, errs.ErrUserAlreadyExists)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(authdto.AuthResponse{Token: token})

}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req authdto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errdto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		errdto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	token, err := h.svc.Login(ctx, &req)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			errdto.WriteJSONError(w, http.StatusUnauthorized, err)
		} else {
			errdto.WriteJSONError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authdto.AuthResponse{Token: token})
}
