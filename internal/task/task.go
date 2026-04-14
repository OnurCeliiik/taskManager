package task

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Status      string         `json:"status" gorm:"not null;default:'pending'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
