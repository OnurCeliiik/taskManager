package task

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateTask(task *Task) (*Task, error)
	GetTaskByID(id uuid.UUID) (*Task, error)
	ListTasksByUserID(userID uuid.UUID, limit, offset int, status, category string) ([]Task, error)
	UpdateTask(task *Task) (*Task, error)
	DeleteTask(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateTask(task *Task) (*Task, error) {
	if err := r.db.Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (r *repository) GetTaskByID(id uuid.UUID) (*Task, error) {
	var task Task
	if err := r.db.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *repository) ListTasksByUserID(userID uuid.UUID, limit, offset int, status, category string) ([]Task, error) {
	var tasks []Task
	query := r.db.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if limit <= 0 {
		limit = 50
	}

	query = query.Limit(limit)

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Order("created_at desc").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *repository) UpdateTask(task *Task) (*Task, error) {
	if err := r.db.Save(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (r *repository) DeleteTask(id uuid.UUID) error {
	return r.db.Delete(&Task{}, "id = ?", id).Error
}
