package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/platonso/hrmate/internal/config"
	"github.com/pressly/goose"
)

var command = flag.String("command", "up", "goose command (up, down, status, etc.)")

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	if _, err := os.Stat(cfg.MigrationsPath); os.IsNotExist(err) {
		log.Fatal("migrations directory does not exist")
	}

	connStr := cfg.GetConnStr()
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("error opening database: ", err)
	}
	defer sqlDB.Close()

	if err = sqlDB.Ping(); err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	if err = goose.Run(*command, sqlDB, cfg.MigrationsPath); err != nil {
		log.Println("failed to run goose command: ", err)
		return
	}

	log.Println("goose command executed successfully")
}
