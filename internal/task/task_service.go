package task

import "github.com/google/uuid"

type TaskService interface {
	CreateTask(req CreateTaskRequest, userID uuid.UUID) (*Task, error)
	GetTaskByID(id uuid.UUID) (*Task, error)
	UpdateTask(id uuid.UUID, req UpdateTaskRequest) (*Task, error)
	DeleteTask(id uuid.UUID) error
}

type taskService struct {
	repository TaskRepository
}

func NewTaskService(repository TaskRepository) TaskService {
	return &taskService{repository: repository}
}

func (s *taskService) GetTaskByID(id uuid.UUID) (*Task, error) {
	return s.repository.GetTaskByID(id)
}

func (s *taskService) CreateTask(req CreateTaskRequest, userID uuid.UUID) (*Task, error) {
	task := &Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Status:      "pending",
	}

	if err := s.repository.CreateTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) UpdateTask(id uuid.UUID, req UpdateTaskRequest) (*Task, error) {
	// Fetch existing task
	task, err := s.repository.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
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
		task.Status = *req.Status
	}

	// Save updates
	if err := s.repository.UpdateTask(id, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) DeleteTask(id uuid.UUID) error {
	// Check if task exists before deletion
	if _, err := s.repository.GetTaskByID(id); err != nil {
		return err
	}
	return s.repository.DeleteTask(id)
}
