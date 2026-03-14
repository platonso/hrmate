package form

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	formdto "github.com/platonso/hrmate/internal/handler/form/dto"
	"github.com/platonso/hrmate/internal/handler/httpapi"
	"github.com/platonso/hrmate/internal/handler/httpapi/dto"
)

type Service interface {
	Create(ctx context.Context, formDTO *formdto.FormCreateRequest, userID uuid.UUID) (*domain.Form, error)
	GetForm(ctx context.Context, formID, currentUserID uuid.UUID) (*formdto.FormWithUserResponse, error)
	GetAllForms(ctx context.Context, currentUserID uuid.UUID) ([]formdto.FormsWithUserResponse, error)
	Update(ctx context.Context, newStatus *formdto.FormStatusUpdateRequest, formID uuid.UUID) (*domain.Form, error)
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

func (h *Handler) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req formdto.FormCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
	}

	userID, ok := httpapi.GetUserID(ctx)
	if !ok {
		dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("missing user ID in context"))
		return
	}

	form, err := h.svc.Create(ctx, &req, userID)
	if err != nil {
		dto.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(form)
}

func (h *Handler) HandleGetForm(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	formIdStr := chi.URLParam(r, "id")
	formID, err := uuid.Parse(formIdStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	userID, _ := httpapi.GetUserID(ctx)

	result, err := h.svc.GetForm(ctx, formID, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrFormNotFound):
			dto.WriteJSONError(w, http.StatusNotFound, errors.New("form not found"))
		case errors.Is(err, errs.ErrForbidden):
			dto.WriteJSONError(w, http.StatusForbidden, errors.New("forbidden"))
		default:
			dto.WriteJSONError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)

}

func (h *Handler) HandleGetForms(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userID, _ := httpapi.GetUserID(ctx)

	result, err := h.svc.GetAllForms(ctx, userID)
	if err != nil {
		dto.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (h *Handler) HandleUpdateFormStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	formIDStr := chi.URLParam(r, "id")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	var req formdto.FormStatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	form, err := h.svc.Update(ctx, &req, formID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrFormNotFound):
			dto.WriteJSONError(w, http.StatusNotFound, errors.New("form not found"))
		default:
			dto.WriteJSONError(w, http.StatusInternalServerError, errors.New("internal server error"))
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(form)

}
