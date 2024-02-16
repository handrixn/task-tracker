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
	GetTaskByUUID(taskUUID string) (*model.Task, error)
	UpdateTask(taskID int64, task *model.Task) (*model.Task, error)
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

func (r *taskRepository) GetTaskByUUID(taskUUID string) (*model.Task, error) {
	var task model.Task
	query := "SELECT id, uuid, title, description, due_date, status, version, created_at, updated_at FROM tasks WHERE uuid=?"
	err := r.db.QueryRow(query, taskUUID).Scan(
		&task.ID,
		&task.UUID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Status,
		&task.Version,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (tr *taskRepository) UpdateTask(taskID int64, task *model.Task) (*model.Task, error) {
	now := time.Now()
	command := "UPDATE tasks SET title=?, description=?, due_date=?, status=?, version=version+1, updated_at=? WHERE id=? AND version=?"
	stmt, err := tr.db.Prepare(command)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(task.Title, task.Description, task.DueDate, task.Status, now, taskID, task.Version)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("no rows affected, update failed")
	}

	task.Version = task.Version + 1
	task.UpdatedAt = now

	return task, nil
}
