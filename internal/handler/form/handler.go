package form

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	errs "github.com/platonso/hrmate/internal/errors"
	formdto "github.com/platonso/hrmate/internal/handler/form/dto"
	"github.com/platonso/hrmate/internal/handler/middleware"
	"github.com/platonso/hrmate/internal/handler/middleware/dto"
	formservice "github.com/platonso/hrmate/internal/service/form"
	"github.com/platonso/hrmate/internal/service/form/model"
)

const (
	QueryParamUserID = "user_id"
	QueryParamStatus = "status"
)

type Service interface {
	Create(ctx context.Context, formDTO *model.FormCreateInput, userID uuid.UUID) (*domain.Form, error)
	GetForm(ctx context.Context, formID uuid.UUID, requesterID uuid.UUID, requesterRole domain.Role) (*domain.Form, error)
	GetForms(ctx context.Context, filter *formservice.Filter, requesterID uuid.UUID, requesterRole domain.Role) ([]domain.Form, error)
	GetFormsWithUsers(ctx context.Context, filter *formservice.Filter, requesterRole domain.Role) ([]model.FormsWithUser, error)
	Approve(ctx context.Context, formID uuid.UUID, comment string) (*domain.Form, error)
	Reject(ctx context.Context, formID uuid.UUID, comment string) (*domain.Form, error)
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

	userID, _ := middleware.GetUserID(r.Context())

	formCreateInput := formdto.ToFormCreateInput(req)

	form, err := h.svc.Create(r.Context(), &formCreateInput, userID)
	if err != nil {
		dto.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(formdto.ToFormResponse(form))
}

func (h *Handler) HandleGetForm(w http.ResponseWriter, r *http.Request) {
	formIdStr := chi.URLParam(r, "id")
	formID, err := uuid.Parse(formIdStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid form id"))
		return
	}

	requesterID, _ := middleware.GetUserID(r.Context())
	requesterRole, _ := middleware.GetUserRole(r.Context())

	form, err := h.svc.GetForm(r.Context(), formID, requesterID, requesterRole)
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
	_ = json.NewEncoder(w).Encode(formdto.ToFormResponse(form))
}

func (h *Handler) HandleGetForms(w http.ResponseWriter, r *http.Request) {
	requesterID, _ := middleware.GetUserID(r.Context())
	requesterRole, _ := middleware.GetUserRole(r.Context())

	// Parse query parameters
	filter, err := h.parseFilter(r)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	// Call service
	forms, err := h.svc.GetForms(r.Context(), filter, requesterID, requesterRole)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrForbidden):
			dto.WriteJSONError(w, http.StatusForbidden, err)
		case errors.Is(err, errs.ErrUserNotFound):
			dto.WriteJSONError(w, http.StatusNotFound, err)
		default:
			dto.WriteJSONError(w, http.StatusInternalServerError, errors.New("internal server error"))
			log.Println(err)
		}
		return
	}

	// Return results (can be empty array)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(formdto.ToFormResponses(forms))
}

func (h *Handler) HandleGetFormsWithUsers(w http.ResponseWriter, r *http.Request) {
	requesterRole, _ := middleware.GetUserRole(r.Context())

	// Parse query parameters
	filter, err := h.parseFilter(r)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	// Call service
	formsWithUsers, err := h.svc.GetFormsWithUsers(r.Context(), filter, requesterRole)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrForbidden):
			dto.WriteJSONError(w, http.StatusForbidden, errors.New("forbidden"))
		case errors.Is(err, errs.ErrUserNotFound):
			dto.WriteJSONError(w, http.StatusNotFound, errors.New("user not found"))
		case strings.Contains(err.Error(), "invalid status"):
			dto.WriteJSONError(w, http.StatusBadRequest, err)
		default:
			dto.WriteJSONError(w, http.StatusInternalServerError, errors.New("internal server error"))
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(formdto.ToFormsWithUserResponses(formsWithUsers))
}

func (h *Handler) HandleApprove(w http.ResponseWriter, r *http.Request) {
	h.handleFormAction(w, r, h.svc.Approve)
}
func (h *Handler) HandleReject(w http.ResponseWriter, r *http.Request) {
	h.handleFormAction(w, r, h.svc.Reject)
}

func (h *Handler) parseFilter(r *http.Request) (*formservice.Filter, error) {
	filter := &formservice.Filter{}

	// Parse user_id
	if userIDStr := r.URL.Query().Get(QueryParamUserID); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid user_id format: %s", userIDStr)
		}
		filter.UserID = &userID
	}

	// Parse status
	if statusStr := r.URL.Query().Get(QueryParamStatus); statusStr != "" {
		status := domain.FormStatus(statusStr)
		filter.FormStatus = &status
		if err := filter.ValidateStatus(); err != nil {
			return nil, fmt.Errorf("invalid status: %s", statusStr)
		}
	}

	return filter, nil
}

func (h *Handler) handleFormAction(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, formID uuid.UUID, comment string) (*domain.Form, error),
) {
	// Parsing formID
	formIDStr := chi.URLParam(r, "id")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}
	// Decoding and validation request
	var req formdto.FormCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}
	if err := h.validator.Struct(req); err != nil {
		dto.WriteJSONError(w, http.StatusBadRequest, errors.New("validation failed"))
		return
	}

	form, err := action(r.Context(), formID, req.Comment)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrFormNotFound):
			dto.WriteJSONError(w, http.StatusNotFound, errs.ErrFormNotFound)
		case errors.Is(err, errs.ErrFormAlreadyApproved):
			dto.WriteJSONError(w, http.StatusBadRequest, errs.ErrFormAlreadyApproved)
		case errors.Is(err, errs.ErrFormAlreadyRejected):
			dto.WriteJSONError(w, http.StatusBadRequest, errs.ErrFormAlreadyRejected)
		default:
			dto.WriteJSONError(w, http.StatusInternalServerError, errors.New("internal server error"))
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(formdto.ToFormResponse(form))
}
