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
	"github.com/platonso/hrmate/internal/handler/middleware"
	"github.com/platonso/hrmate/internal/handler/middleware/dto"
	"github.com/platonso/hrmate/internal/service/form/model"
)

type Service interface {
	Create(ctx context.Context, formDTO *model.FormCreateInput, userID uuid.UUID) (*domain.Form, error)
	GetForm(ctx context.Context, formID, currentUserID uuid.UUID) (*domain.Form, error)
	//	GetFormsWithUser(ctx context.Context, targetUserID, requesterUserID uuid.UUID) ([]model.FormsWithUser, error)
	GetFormsWithUsers(ctx context.Context, currentUserID uuid.UUID) ([]model.FormsWithUser, error)
	Approve(ctx context.Context, formID uuid.UUID, comment string) (*domain.Form, error)
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
	var req formdto.FormCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		dto.WriteJSONError(w, http.StatusUnauthorized, errors.New("missing user ID in context"))
		return
	}

	form, err := h.svc.Create(r.Context(), formdto.ToFormCreateInput(&req), userID)
	if err != nil {
		dto.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(form)
}

func (h *Handler) HandleGetForm(w http.ResponseWriter, r *http.Request) {
	formIdStr := chi.URLParam(r, "id")
	formID, err := uuid.Parse(formIdStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	userID, _ := middleware.GetUserID(r.Context())

	form, err := h.svc.GetForm(r.Context(), formID, userID)
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
	_ = json.NewEncoder(w).Encode(form)
}

func (h *Handler) HandleGetForms(w http.ResponseWriter, r *http.Request) {
	h.svc.GetForms()
}

func (h *Handler) HandleGetFormsWithUsers(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	forms, err := h.svc.GetFormsWithUsers(r.Context(), userID)
	if err != nil {
		dto.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(forms)
}

func (h *Handler) HandleApprove(w http.ResponseWriter, r *http.Request) {
	formIDStr := chi.URLParam(r, "id")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	var req formdto.FormCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	form, err := h.svc.Approve(r.Context(), formID, req.Comment)
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
