package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/platonso/hrmate/internal/config"
	"github.com/pressly/goose/v3"
)

var command = flag.String("command", "up", "goose command (up, down, status)")

func main() {
	flag.Parse()

	cfg, err := config.NewDB()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	if _, err := os.Stat(cfg.MigrationDir); os.IsNotExist(err) {
		log.Fatal("migrations directory does not exist")
	}

	connStr := cfg.GetConnStr()
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("error opening database: ", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	if err = sqlDB.Ping(); err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	switch *command {
	case "up":
		if err := goose.Up(sqlDB, cfg.MigrationDir); err != nil {
			log.Fatalf("failed to run up: %v", err)
		}
	case "down":
		if err := goose.Down(sqlDB, cfg.MigrationDir); err != nil {
			log.Fatalf("failed to run down: %v", err)
		}
	case "status":
		if err := goose.Status(sqlDB, cfg.MigrationDir); err != nil {
			log.Fatalf("failed to run status: %v", err)
		}
	default:
		log.Fatalf("unknown command: %s", *command)
	}
}
