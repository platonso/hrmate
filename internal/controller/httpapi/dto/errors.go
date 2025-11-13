package dto

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ErrResponseDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func WriteJSONError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResponseDTO := ErrResponseDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	if err := json.NewEncoder(w).Encode(errResponseDTO); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}
