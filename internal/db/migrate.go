package db

import (
	"database/sql"
	"log"
)

func RunMigrations(conn *sql.DB) {
	usersSchema := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		user_role TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		position TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		hashed_password TEXT NOT NULL,
        is_active BOOLEAN NOT NULL
	);`

	formsSchema := `
	CREATE TABLE IF NOT EXISTS forms (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		start_date TIMESTAMPTZ,
		end_date TIMESTAMPTZ,
		created_at TIMESTAMPTZ NOT NULL,
		approved_at TIMESTAMPTZ,
		status TEXT NOT NULL
	);`

	if _, err := conn.Exec(usersSchema); err != nil {
		log.Fatal("failed to create users table:", err)
	}

	if _, err := conn.Exec(formsSchema); err != nil {
		log.Fatal("failed to create forms table:", err)
	}
}
