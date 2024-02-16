package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/handrixn/task-tracker/internal/constant"
	"github.com/handrixn/task-tracker/internal/model"
)

type TaskRepository interface {
	Create(m *model.Task) (*model.Task, error)
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *taskRepository {
	return &taskRepository{db: db}
}

func (tr *taskRepository) Create(task *model.Task) (*model.Task, error) {
	command := `
		INSERT INTO tasks (uuid, title, description, due_date, status, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := tr.db.Prepare(command)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	generatedUUID := uuid.New().String()
	now := time.Now()

	result, err := stmt.Exec(generatedUUID, task.Title, task.Description, task.DueDate, task.Status, 1, now, now)

	if err != nil {
		log.Println(err)
		return nil, errors.New(constant.FAILED_INSERT_DATA)
	}

	task.ID, _ = result.LastInsertId()
	task.UUID = generatedUUID
	task.CreatedAt = now
	task.UpdatedAt = now
	task.Version = 1

	return task, nil
}
