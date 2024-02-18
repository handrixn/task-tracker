package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/handrixn/task-tracker/internal/model"
	"github.com/handrixn/task-tracker/internal/repository"
)

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	now := time.Now()

	command := `INSERT INTO tasks`
	mock.ExpectPrepare(command).
		ExpectExec().
		WithArgs(sqlmock.AnyArg(), "Test Task Title", "Test Task Description", now, "in-progress", 1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	taskRepo := repository.NewTaskRepository(db)

	task := &model.Task{
		Title:       "Test Task Title",
		Description: "Test Task Description",
		DueDate:     now,
		Status:      "in-progress",
		Version:     1,
	}

	createdTask, err := taskRepo.Create(task)
	if err != nil {
		t.Errorf("Error creating task: %v", err)
	}

	if createdTask == nil {
		t.Error("Expected created task, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "uuid", "title", "description", "due_date", "status", "version", "created_at", "updated_at"}).
		AddRow(1, "test-uuid", "Test Task Title", "Test Task Description", now, "in-progress", 1, now, now)

	mock.ExpectQuery("SELECT").WithArgs("test-uuid").WillReturnRows(rows)

	taskRepo := repository.NewTaskRepository(db)

	task, err := taskRepo.GetByUUID("test-uuid")
	if err != nil {
		t.Errorf("Error getting task by UUID: %v", err)
	}

	if task == nil {
		t.Error("Expected task, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestUpdateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	now := time.Now()

	mock.
		ExpectPrepare("UPDATE tasks").
		ExpectExec().
		WithArgs("Updated Task Title", "Updated Task Description", now, "completed", sqlmock.AnyArg(), 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	taskRepo := repository.NewTaskRepository(db)

	task := &model.Task{
		ID:          1,
		Title:       "Updated Task Title",
		Description: "Updated Task Description",
		DueDate:     now,
		Status:      "completed",
		Version:     1,
	}

	updatedTask, err := taskRepo.Update(1, task)
	if err != nil {
		t.Errorf("Error updating task: %v", err)
	}

	if updatedTask == nil {
		t.Error("Expected updated task, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestListTasks(t *testing.T) {
	tests := []struct {
		name              string
		rows              *sqlmock.Rows
		filter            map[string]string
		expectedTasks     int
		expectedQueryArgs []interface{}
	}{
		{
			name: "No filter",
			rows: sqlmock.NewRows([]string{"id", "uuid", "title", "description", "due_date", "status"}).
				AddRow(1, "test-uuid-1", "Test Task 1", "Test Description 1", time.Now(), "in-progress").
				AddRow(2, "test-uuid-2", "Test Task 2", "Test Description 2", time.Now(), "completed"),
			expectedTasks:     2,
			expectedQueryArgs: nil,
		},
		{
			name: "Filter by title",
			rows: sqlmock.NewRows([]string{"id", "uuid", "title", "description", "due_date", "status"}).
				AddRow(1, "test-uuid-1", "Test Task 1", "Test Description 1", time.Now(), "in-progress"),
			filter:            map[string]string{"title": "Test Task 1"},
			expectedTasks:     1,
			expectedQueryArgs: []interface{}{"%Test Task 1%"},
		},
		{
			name: "Filter by status",
			rows: sqlmock.NewRows([]string{"id", "uuid", "title", "description", "due_date", "status"}).
				AddRow(2, "test-uuid-2", "Test Task 2", "Test Description 2", time.Now(), "completed"),
			filter:            map[string]string{"status": "completed"},
			expectedTasks:     1,
			expectedQueryArgs: []interface{}{"completed"},
		},
		{
			name: "Filter by due date",
			rows: sqlmock.NewRows([]string{"id", "uuid", "title", "description", "due_date", "status"}).
				AddRow(1, "test-uuid-1", "Test Task 1", "Test Description 1", time.Now(), "in-progress"),
			filter:            map[string]string{"due_date": "2024-02-20"},
			expectedTasks:     1,
			expectedQueryArgs: []interface{}{"2024-02-20"},
		},
		{
			name: "Multiple filters",
			rows: sqlmock.NewRows([]string{"id", "uuid", "title", "description", "due_date", "status"}).
				AddRow(1, "test-uuid-1", "Test Task 1", "Test Description 1", time.Now(), "in-progress"),
			filter:            map[string]string{"title": "Test", "status": "in-progress"},
			expectedTasks:     1,
			expectedQueryArgs: []interface{}{"%Test%", "in-progress"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error creating mock database: %v", err)
			}
			defer db.Close()

			if tt.rows != nil {
				mock.ExpectQuery("SELECT").WillReturnRows(tt.rows)
			}

			taskRepo := repository.NewTaskRepository(db)

			tasks, err := taskRepo.List(tt.filter, 0, 0)
			if err != nil {
				t.Errorf("Error listing tasks: %v", err)
			}

			if len(tasks) != tt.expectedTasks {
				t.Errorf("Expected %d tasks, got %d", tt.expectedTasks, len(tasks))
			}

			if tt.expectedQueryArgs != nil {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("Unfulfilled expectations: %s", err)
				}
			}
		})
	}
}

func TestCountTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"COUNT(*)"}).
		AddRow(2)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM tasks").WillReturnRows(rows)

	taskRepo := repository.NewTaskRepository(db)

	count, err := taskRepo.Count(nil)

	if err != nil {
		t.Errorf("Error counting tasks: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetTaskSummary(t *testing.T) {
	tests := []struct {
		name               string
		rows               *sqlmock.Rows
		filter             map[string]string
		expectedTotal      int
		expectedInProgress int
		expectedCompleted  int
		expectedQueryArgs  []interface{}
	}{
		{
			name: "No filter",
			rows: sqlmock.NewRows([]string{"total", "in_progress", "completed"}).
				AddRow(5, 3, 2),
			expectedTotal:      5,
			expectedInProgress: 3,
			expectedCompleted:  2,
			expectedQueryArgs:  nil,
		},
		{
			name: "Filter by due date",
			rows: sqlmock.NewRows([]string{"total", "in_progress", "completed"}).
				AddRow(3, 2, 1),
			filter:             map[string]string{"due_date": "2024-02-20"},
			expectedTotal:      3,
			expectedInProgress: 2,
			expectedCompleted:  1,
			expectedQueryArgs:  []interface{}{"2024-02-20"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error creating mock database: %v", err)
			}
			defer db.Close()

			if tt.rows != nil {
				mock.ExpectQuery("SELECT").WillReturnRows(tt.rows)
			}

			taskRepo := repository.NewTaskRepository(db)

			summary, err := taskRepo.GetSummary(tt.filter)
			if err != nil {
				t.Errorf("Error getting task summary: %v", err)
			}

			if summary.Total != tt.expectedTotal || summary.InProgress != tt.expectedInProgress || summary.Completed != tt.expectedCompleted {
				t.Errorf("Unexpected task summary values. Expected total=%d, in_progress=%d, completed=%d. Got total=%d, in_progress=%d, completed=%d", tt.expectedTotal, tt.expectedInProgress, tt.expectedCompleted, summary.Total, summary.InProgress, summary.Completed)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
