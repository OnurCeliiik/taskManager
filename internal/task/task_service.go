package task

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskService interface {
	CreateTask(req CreateTaskRequest, userID uuid.UUID) (*TaskResponse, error)
	ListTasks(userID uuid.UUID, limit, offset int, status, category string) ([]TaskResponse, error)
	GetTaskByID(id uuid.UUID) (*TaskResponse, error)
	UpdateTask(taskID uuid.UUID, userID uuid.UUID, req UpdateTaskRequest) (*TaskResponse, error)
	DeleteTask(taskID uuid.UUID, userID uuid.UUID) error
}

type taskService struct {
	repository Repository
}

func NewTaskService(repository Repository) TaskService {
	return &taskService{repository: repository}
}

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskUnauthorized  = errors.New("task not owned by user")
	ErrInvalidTaskStatus = errors.New("invalid task status")
)

func (s *taskService) CreateTask(req CreateTaskRequest, userID uuid.UUID) (*TaskResponse, error) {
	task := &Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Status:      "pending",
	}

	created, err := s.repository.CreateTask(task)
	if err != nil {
		return nil, err
	}

	return toTaskResponse(created), nil
}

func (s *taskService) ListTasks(userID uuid.UUID, limit, offset int, status, category string) ([]TaskResponse, error) {
	tasks, err := s.repository.ListTasksByUserID(userID, limit, offset, status, category)
	if err != nil {
		return nil, err
	}

	response := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, *toTaskResponse(&task))
	}
	return response, nil
}

func (s *taskService) GetTaskByID(id uuid.UUID) (*TaskResponse, error) {
	task, err := s.repository.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return toTaskResponse(task), nil
}

func (s *taskService) UpdateTask(taskID uuid.UUID, userID uuid.UUID, req UpdateTaskRequest) (*TaskResponse, error) {
	task, err := s.repository.GetTaskByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if task.UserID != userID {
		return nil, ErrTaskUnauthorized
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Category != nil {
		task.Category = *req.Category
	}
	if req.Status != nil {
		status := *req.Status
		if status != "pending" && status != "in_progress" && status != "completed" {
			return nil, ErrInvalidTaskStatus
		}
		task.Status = status
	}

	task.UpdatedAt = time.Now()

	updated, err := s.repository.UpdateTask(task)
	if err != nil {
		return nil, err
	}

	return toTaskResponse(updated), nil
}

func (s *taskService) DeleteTask(taskID uuid.UUID, userID uuid.UUID) error {
	task, err := s.repository.GetTaskByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	if task.UserID != userID {
		return ErrTaskUnauthorized
	}

	return s.repository.DeleteTask(taskID)
}

func toTaskResponse(task *Task) *TaskResponse {
	return &TaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Category:    task.Category,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}
}
