package task

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task *Task) error
	GetTaskByID(id uuid.UUID) (*Task, error)
	UpdateTask(id uuid.UUID, task *Task) error
	DeleteTask(id uuid.UUID) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task *Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) GetTaskByID(id uuid.UUID) (*Task, error) {
	var task Task
	if err := r.db.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) UpdateTask(id uuid.UUID, task *Task) error {
	return r.db.Model(&Task{}).Where("id = ?", id).Updates(task).Error
}

func (r *taskRepository) DeleteTask(id uuid.UUID) error {
	return r.db.Delete(&Task{}, "id = ?", id).Error
}
