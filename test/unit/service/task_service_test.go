package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/handrixn/task-tracker/internal/constant"
	"github.com/handrixn/task-tracker/internal/model"
	"github.com/handrixn/task-tracker/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task *model.Task) (*model.Task, error) {
	args := m.Called(task)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByUUID(taskUUID string) (*model.Task, error) {
	args := m.Called(taskUUID)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(taskID int64, task *model.Task) (*model.Task, error) {
	args := m.Called(taskID, task)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) List(filter map[string]string, page, limit int) ([]model.Task, error) {
	args := m.Called(filter, page, limit)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]model.Task), args.Error(1)
}

func (m *MockTaskRepository) Count(filter map[string]string) (int64, error) {
	args := m.Called(filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTaskRepository) GetSummary(filter map[string]string) (*model.TaskSummary, error) {
	args := m.Called(filter)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*model.TaskSummary), args.Error(1)
}

func TestCreateTask(t *testing.T) {
	testCases := []struct {
		name            string
		inputTask       *model.TaskInput
		expectedTask    *model.Task
		expectedErr     error
		repositoryError error
	}{
		{
			name: "Success",
			inputTask: &model.TaskInput{
				Title:       "Test Title",
				Description: "Test Description",
				DueDate:     "2024-02-20",
			},
			expectedTask: &model.Task{
				ID:          1,
				Title:       "Test Title",
				Description: "Test Description",
				DueDate:     time.Time{},
				Status:      constant.IN_PROGRESS,
			},
			expectedErr: nil,
		},
		{
			name: "Repository Error",
			inputTask: &model.TaskInput{
				Title:       "Test Title",
				Description: "Test Description",
				DueDate:     "2024-02-20",
			},
			expectedTask:    nil,
			expectedErr:     errors.New("repository error"),
			repositoryError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			taskService := service.NewTaskService(mockRepo)

			mockRepo.On("Create", mock.Anything).Return(tc.expectedTask, tc.repositoryError)

			createdTask, err := taskService.CreateTask(tc.inputTask)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedTask, createdTask)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	testCases := []struct {
		name            string
		taskUUID        string
		updateTaskInput *model.TaskInputUpdate
		expectedTask    *model.Task
		expectedErr     error
		repositoryError error
	}{
		{
			name:     "Success",
			taskUUID: "test-uuid",
			updateTaskInput: &model.TaskInputUpdate{
				TaskInput: model.TaskInput{
					Title:       "Updated Title",
					Description: "Updated Description",
					DueDate:     "2024-02-20",
				},
				Status: constant.COMPLETED,
			},
			expectedTask: &model.Task{
				ID:          1,
				Title:       "Updated Title",
				Description: "Updated Description",
				DueDate:     time.Time{},
				Status:      constant.COMPLETED,
			},
			expectedErr: nil,
		},
		{
			name:     "GetByUUID Error",
			taskUUID: "test-uuid",
			updateTaskInput: &model.TaskInputUpdate{
				TaskInput: model.TaskInput{
					Title:       "Updated Title",
					Description: "Updated Description",
					DueDate:     "2024-02-20",
				},
				Status: constant.COMPLETED,
			},
			expectedTask:    nil,
			expectedErr:     errors.New("repository error"),
			repositoryError: errors.New("repository error"),
		},
		{
			name:     "Update Error",
			taskUUID: "test-uuid",
			updateTaskInput: &model.TaskInputUpdate{
				TaskInput: model.TaskInput{
					Title:       "Updated Title",
					Description: "Updated Description",
					DueDate:     "2024-02-20",
				},
				Status: constant.COMPLETED,
			},
			expectedTask:    nil,
			expectedErr:     errors.New("repository error"),
			repositoryError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			taskService := service.NewTaskService(mockRepo)

			mockRepo.On("GetByUUID", tc.taskUUID).Return(&model.Task{ID: 1}, tc.repositoryError)
			mockRepo.On("Update", mock.Anything, mock.Anything).Return(tc.expectedTask, tc.repositoryError)

			updatedTask, err := taskService.UpdateTask(tc.taskUUID, tc.updateTaskInput)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedTask, updatedTask)
		})
	}
}

func TestListTasks(t *testing.T) {
	testCases := []struct {
		name          string
		params        map[string]string
		page          int
		limit         int
		expectedTasks []model.Task
		expectedCount int64
		expectedErr   error
		repositoryErr error
	}{
		{
			name: "Success",
			params: map[string]string{
				"page":  "1",
				"limit": "10",
			},
			page:  1,
			limit: 10,
			expectedTasks: []model.Task{
				{ID: 1, Title: "Task 1", Description: "Description 1", DueDate: time.Now(), Status: constant.IN_PROGRESS},
				{ID: 2, Title: "Task 2", Description: "Description 2", DueDate: time.Now(), Status: constant.COMPLETED},
			},
			expectedCount: 2,
			expectedErr:   nil,
		},
		{
			name: "Error",
			params: map[string]string{
				"page":  "1",
				"limit": "10",
			},
			page:          1,
			limit:         10,
			expectedTasks: nil,
			expectedCount: 0,
			expectedErr:   errors.New("repository error"),
			repositoryErr: errors.New("repository error"),
		},
		{
			name: "SearchByTitle",
			params: map[string]string{
				"page":  "1",
				"limit": "10",
				"title": "Task",
			},
			page:  1,
			limit: 10,
			expectedTasks: []model.Task{
				{ID: 1, Title: "Task 1", Description: "Description 1", DueDate: time.Now(), Status: constant.IN_PROGRESS},
				{ID: 2, Title: "Task 2", Description: "Description 2", DueDate: time.Now(), Status: constant.COMPLETED},
			},
			expectedCount: 2,
			expectedErr:   nil,
		},
		{
			name: "FilterByStatus",
			params: map[string]string{
				"page":   "1",
				"limit":  "10",
				"status": constant.IN_PROGRESS,
			},
			page:  1,
			limit: 10,
			expectedTasks: []model.Task{
				{ID: 1, Title: "Task 1", Description: "Description 1", DueDate: time.Now(), Status: constant.IN_PROGRESS},
			},
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name: "FilterByDueDate",
			params: map[string]string{
				"page":     "1",
				"limit":    "10",
				"due_date": "2024-02-20",
			},
			page:  1,
			limit: 10,
			expectedTasks: []model.Task{
				{ID: 1, Title: "Task 1", Description: "Description 1", DueDate: time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC), Status: constant.IN_PROGRESS},
				{ID: 2, Title: "Task 2", Description: "Description 2", DueDate: time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC), Status: constant.COMPLETED},
			},
			expectedCount: 2,
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			taskService := service.NewTaskService(mockRepo)

			mockRepo.On("List", tc.params, tc.page, tc.limit).Return(tc.expectedTasks, tc.repositoryErr)
			mockRepo.On("Count", tc.params).Return(tc.expectedCount, tc.repositoryErr)

			tasks, totalCount, err := taskService.ListTasks(tc.params)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedTasks, tasks)
			assert.Equal(t, tc.expectedCount, totalCount)
		})
	}
}

func TestTaskSummary(t *testing.T) {
	testCases := []struct {
		name            string
		params          map[string]string
		expectedSummary *model.TaskSummary
		expectedErr     error
	}{
		{
			name: "ValidParams",
			params: map[string]string{
				"due_date": "2024-02-20",
			},
			expectedSummary: &model.TaskSummary{
				Total:      5,
				InProgress: 3,
				Completed:  2,
			},
			expectedErr: nil,
		},
		{
			name:            "EmptyParams",
			params:          map[string]string{},
			expectedSummary: nil,
			expectedErr:     errors.New("no parameters provided"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskRepo := new(MockTaskRepository)
			taskService := service.NewTaskService(taskRepo)

			taskRepo.On("GetSummary", tc.params).Return(tc.expectedSummary, tc.expectedErr)

			summary, err := taskService.TaskSummary(tc.params)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedSummary, summary)
		})
	}
}
