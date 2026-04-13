package users

import "github.com/google/uuid"

type Service interface {
	GetUserByEmail(email string) (*User, error)
	// UpdateUser(id string, req UpdateUserRequest) (*UserResponse, error)
	GetUserByID(id uuid.UUID) (*User, error)
	// DeleteUser(id string) error
}

type service struct {
	repository Repository
}

func NewUserService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) GetUserByEmail(email string) (*User, error) {
	return s.repository.GetUserByEmail(email)
}

func (s *service) GetUserByID(id uuid.UUID) (*User, error) {
	return s.repository.GetUserByID(id)
}
