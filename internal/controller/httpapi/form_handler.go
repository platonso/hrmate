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

type FormHandler struct {
	formService *service.FormService
	validator   *validator.Validate
}

func NewFormHandler(formService *service.FormService, v *validator.Validate) *FormHandler {
	return &FormHandler{
		formService: formService,
		validator:   v,
	}
}

func (h *FormHandler) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req dto.FormCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
	}

	userID, ok := GetUserID(ctx)
	if !ok {
		dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("missing user ID in context"))
		return
	}

	form, err := h.formService.Create(ctx, &req, userID)
	if err != nil {
		dto.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(form)
}

func (h *FormHandler) HandleGetForm(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	formIdStr := chi.URLParam(r, "id")
	formID, err := uuid.Parse(formIdStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	userID, _ := GetUserID(ctx)

	result, err := h.formService.GetForm(ctx, formID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFormNotFound):
			dto.WriteJSONError(w, http.StatusNotFound, nil)
		case errors.Is(err, domain.ErrForbidden):
			dto.WriteJSONError(w, http.StatusForbidden, err)
		default:
			dto.WriteJSONError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)

}

func (h *FormHandler) HandleGetForms(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userID, _ := GetUserID(ctx)

	result, err := h.formService.GetAllForms(ctx, userID)
	if err != nil {
		dto.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (h *FormHandler) HandleUpdateFormStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	formIDStr := chi.URLParam(r, "id")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	var req dto.FormStatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	userID, _ := GetUserID(ctx)

	form, err := h.formService.Update(ctx, &req, formID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFormNotFound):
			dto.WriteJSONError(w, http.StatusNotFound, nil)
		default:
			dto.WriteJSONError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(form)

}
