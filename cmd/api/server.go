package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *application) StartServer() error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.port),
		Handler: app.routes(),
	}

	log.Printf("Starting server on port %d", app.port)
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
