package repository

import (
	"database/sql"
	"errors"
	"log"
	"strings"
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
	List(filter map[string]string, page, limit int) ([]model.Task, error)
	Count(filter map[string]string) (int64, error)
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *taskRepository {
	return &taskRepository{db: db}
}

func getAndBuildListFilter(filter map[string]string, builder *strings.Builder) []any {
	var filterValues []any

	for k, v := range filter {
		if k == constant.TASK_FILTER_TITILE_NAME && v != "" {
			builder.WriteString(" AND title LIKE ?")
			filterValues = append(filterValues, "%"+v+"%")
		}

		if k == constant.TASK_FILTER_STATUS_NAME && v != "" {
			builder.WriteString(" AND status = ?")
			filterValues = append(filterValues, v)
		}

		if k == constant.TASK_FILTER_DUE_DATE_NAME && v != "" {
			builder.WriteString(" AND due_date = ?")
			filterValues = append(filterValues, v)
		}
	}

	return filterValues
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

func (tr *taskRepository) List(filter map[string]string, page, limit int) ([]model.Task, error) {
	var builder strings.Builder

	baseQuery := "SELECT id, uuid, title, description, due_date, status FROM tasks WHERE 1=1"
	builder.WriteString(baseQuery)

	queryParameters := getAndBuildListFilter(filter, &builder)

	if page != 0 && limit != 0 {
		offset := (page - 1) * limit
		queryParameters = append(queryParameters, limit, offset)
		builder.WriteString(" LIMIT ? OFFSET ?")
	}

	query := builder.String()

	rows, err := tr.db.Query(query, queryParameters...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.UUID, &task.Title, &task.Description, &task.DueDate, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return tasks, nil
}

func (tr *taskRepository) Count(filter map[string]string) (int64, error) {
	var builder strings.Builder
	builder.WriteString("SELECT COUNT(*) FROM tasks WHERE 1=1")

	queryParameters := getAndBuildListFilter(filter, &builder)
	countQuery := builder.String()

	var totalCount int64
	err := tr.db.QueryRow(countQuery, queryParameters...).Scan(&totalCount)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return totalCount, err
}
