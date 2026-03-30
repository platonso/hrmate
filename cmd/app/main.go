package main

import (
	"context"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/platonso/hrmate/internal/app"
	"github.com/platonso/hrmate/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	a, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}
	defer a.Close()

	if err := a.StartServer(); err != nil {
		log.Fatalf("start server error: %v", err)
	}

}
