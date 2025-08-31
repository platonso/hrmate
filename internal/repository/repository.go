package repository

import (
	"database/sql"
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/repository/postgres"
)

type Repository struct {
	Users domain.EmployeeRepository
	Forms domain.FormRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Users: &postgres.UserRepository{DB: db},
		Forms: &postgres.FormRepository{DB: db},
	}
}
