package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	"log"
	"net/http"
)

func (app *application) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
	var formDTO FormCreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&formDTO); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	userId, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "missing user ID in context", http.StatusUnauthorized)
		return
	}

	form := domain.NewForm(
		userId,
		formDTO.Title,
		formDTO.Description,
		formDTO.StartDate,
		formDTO.EndDate,
	)

	if err := app.repositories.Forms.Create(&form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("failed to save form: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(form); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (app *application) HandleGetForm(w http.ResponseWriter, r *http.Request) {
	formIdStr := chi.URLParam(r, "id")

	formId, err := uuid.Parse(formIdStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	role, ok := GetUserRole(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusUnauthorized, errors.New("missing role"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch role {
	case domain.RoleAdmin, domain.RoleHR:
		formWithUser, err := app.repositories.Forms.FindByIDWithUser(formId)
		if err != nil {
			if errors.Is(err, domain.ErrFormNotFound) {
				WriteJSONError(w, http.StatusNotFound, domain.ErrFormNotFound)
			} else {
				WriteJSONError(w, http.StatusInternalServerError, err)
			}
			return
		}

		if err := json.NewEncoder(w).Encode(formWithUser); err != nil {
			log.Printf("failed to encode response: %v", err)
		}

	case domain.RoleEmployee:
		userID, ok := GetUserID(r.Context())
		if !ok {
			WriteJSONError(w, http.StatusUnauthorized, errors.New("missing user ID"))
			return
		}

		form, err := app.repositories.Forms.FindByID(formId)
		if err != nil {
			if errors.Is(err, domain.ErrFormNotFound) {
				WriteJSONError(w, http.StatusNotFound, domain.ErrFormNotFound)
			} else {
				WriteJSONError(w, http.StatusInternalServerError, err)
			}
			return
		}
		if userID != form.UserID {
			WriteJSONError(w, http.StatusForbidden, errors.New("access denied"))
			return
		}

		if err := json.NewEncoder(w).Encode(form); err != nil {
			log.Printf("failed to encode response: %v", err)
		}

	default:
		WriteJSONError(w, http.StatusForbidden, errors.New("role not authorized"))
	}
}

func (app *application) HandleGetForms(w http.ResponseWriter, r *http.Request) {
	role, ok := GetUserRole(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusUnauthorized, errors.New("missing role"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch role {
	case domain.RoleAdmin, domain.RoleHR:
		formsWithUsers, err := app.repositories.Forms.FindAllWithUsers()
		if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}

		if len(formsWithUsers) == 0 {
			formsWithUsers = []domain.UserWithForms{}
		}

		if err := json.NewEncoder(w).Encode(formsWithUsers); err != nil {
			log.Printf("failed to encode response: %v", err)
			return
		}

	case domain.RoleEmployee:
		userID, ok := GetUserID(r.Context())
		if !ok {
			WriteJSONError(w, http.StatusUnauthorized, errors.New("missing user ID"))
			return
		}

		formsWithUser, err := app.repositories.Forms.FindByUserIDWithUser(userID)
		if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}

		var forms []domain.Form
		if len(formsWithUser) > 0 {
			forms = formsWithUser[0].Forms
		} else {
			forms = []domain.Form{}
		}

		if err := json.NewEncoder(w).Encode(forms); err != nil {
			log.Printf("failed to encode response: %v", err)
		}

	default:
		WriteJSONError(w, http.StatusForbidden, errors.New("role not authorized"))
		return
	}
}

func (app *application) HandleUpdateFormStatus(w http.ResponseWriter, r *http.Request) {
	formIdStr := chi.URLParam(r, "id")
	formId, err := uuid.Parse(formIdStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	form, err := app.repositories.Forms.FindByID(formId)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, domain.ErrFormNotFound)
		return
	}

	var statusFormDTO FormStatusUpdateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&statusFormDTO); err != nil {
		WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if err := statusFormDTO.Validate(); err != nil {
		WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	form.UpdateStatus(statusFormDTO.Status)

	if err := app.repositories.Forms.Update(form); err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err)
		log.Printf("failed to update form: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(form); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
