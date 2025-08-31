package main

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/platonso/hrmate/internal/db"
	"github.com/platonso/hrmate/internal/env"
	"github.com/platonso/hrmate/internal/repository"
	"log"
	"os"
)

type application struct {
	port         int
	jwtSecret    string
	repositories *repository.Repository
	validator    *validator.Validate
}

func main() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	db.RunMigrations(conn)

	app := &application{
		port:         env.GetEnvInt("PORT", 8080),
		jwtSecret:    os.Getenv("JWT_SECRET"),
		repositories: repository.NewRepository(conn),
		validator:    validator.New(),
	}

	if err := app.adminImplementation(); err != nil {
		log.Fatalf("failed to implement admin: %v", err)
	}

	if err := app.StartServer(); err != nil {
		log.Fatal(err)
	}
}
