package users

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *User) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
	UpdateUser(user *User) (*User, error)
	DeleteUser(ID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user *User) (*User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByID(id uuid.UUID) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUser(user *User) (*User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) DeleteUser(ID uuid.UUID) error {
	return r.db.Delete(&User{}, "id = ?", ID).Error
}
