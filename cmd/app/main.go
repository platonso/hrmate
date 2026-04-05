package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/platonso/hrmate/internal/app"
	"github.com/platonso/hrmate/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Printf("Config error: %v", err)
		os.Exit(1)
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := app.New(signalCtx, cfg)
	if err != nil {
		log.Printf("Failed to init app: %v", err)
		os.Exit(1)
	}

	errChan := make(chan error, 1)
	go a.Start(errChan)

	select {
	case err := <-errChan:
		log.Printf("Server stopped with error: %v", err)
		os.Exit(1)

	case <-signalCtx.Done():
		log.Println("Shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(signalCtx, 10*time.Second)
		defer cancel()

		if err := a.Stop(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			os.Exit(1)
		}
	}
}
