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

	//GS:
	//func main() {
	//	// ... инициализация ...
	//
	//	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	//	defer stop()
	//
	//	errChan := make(chan error, 1)
	//	go a.Start(ctx, errChan)
	//
	//	select {
	//	case err := <-errChan:
	//		logger.Error("failed to start", slog.Any("error", err))
	//		os.Exit(1)
	//	case <-ctx.Done():
	//		logger.Info("shutting down")
	//
	//		// Таймаут для остановки
	//		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	//		defer cancel()
	//
	//		if err := a.Stop(shutdownCtx); err != nil {
	//			logger.Error("shutdown error", slog.Any("error", err))
	//			os.Exit(1)
	//		}
	//	}
	//}

}
