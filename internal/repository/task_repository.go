package repository

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type TaskRepository interface {
	// Implement repository methods
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *taskRepository {
	return &taskRepository{db: db}
}
