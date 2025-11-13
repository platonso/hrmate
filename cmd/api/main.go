package main

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/platonso/hrmate/internal/app"
	"github.com/platonso/hrmate/internal/config"
	"log"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err)
	}
	defer application.Close()

	if err := application.StartServer(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//cfg, err := config.New()
	//if err != nil {
	//	log.Fatalf("Config error: %v", err)
	//}
	//
	//application, err := app.New(ctx, cfg)
	//if err != nil {
	//	log.Fatalf("Failed to init app: %v", err)
	//}
	//defer application.Close()

}
