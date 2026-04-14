package users

import (
	"task-manager/internal/task"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name             string         `json:"name" gorm:"not null"`
	Email            string         `json:"email" gorm:"unique;not null"`
	Password         string         `json:"password" gorm:"not null"`
	Tasks            *[]task.Task   `json:"tasks,omitempty"`
	Role             string         `json:"role" gorm:"not null;default:'user'"`
	ResetToken       *string        `json:"-" gorm:"index"`
	ResetTokenExpiry *time.Time     `json:"-"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
