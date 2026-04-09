package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	errs "github.com/platonso/hrmate/internal/errors"
)

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

func WriteError(w http.ResponseWriter, err error, msg string) {
	if err == nil {
		err = errs.ErrInternalServer
		msg = "unknown error"
	}

	var statusCode int

	switch {
	case errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrUnauthorized):
		statusCode = http.StatusUnauthorized

	case errors.Is(err, errs.ErrUserNotActive),
		errors.Is(err, errs.ErrForbidden):
		statusCode = http.StatusForbidden

	case errors.Is(err, errs.ErrUserAlreadyExists),
		errors.Is(err, errs.ErrFormAlreadyApproved),
		errors.Is(err, errs.ErrFormAlreadyRejected),
		errors.Is(err, errs.ErrNoAvailableExecutors):
		statusCode = http.StatusConflict

	case errors.Is(err, errs.ErrUserNotFound),
		errors.Is(err, errs.ErrFormNotFound):
		statusCode = http.StatusNotFound

	case errors.Is(err, errs.ErrInvalidRequest):
		statusCode = http.StatusBadRequest

	default:
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResponse := errorResponse{
		Error: errorDetail{
			Code:    err.Error(),
			Message: msg,
		},
	}

	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}
