package users

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetUserByEmail(email string) (*User, error)
	UpdateUser(ID uuid.UUID, req UpdateUserRequest) (*UserResponse, error)
	GetUserByID(id uuid.UUID) (*User, error)
	DeleteUser(ID uuid.UUID) error
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

func (s *service) UpdateUser(ID uuid.UUID, req UpdateUserRequest) (*UserResponse, error) {
	existingUser, err := s.repository.GetUserByID(ID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existingUser.Name = *req.Name
	}
	if req.Email != nil {
		existingUser.Email = *req.Email
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("could not hash the password")
		}
		existingUser.Password = string(hashedPassword)
	}
	existingUser.UpdatedAt = time.Now()

	respUser, err := s.repository.UpdateUser(existingUser)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:    respUser.ID,
		Name:  respUser.Name,
		Email: respUser.Email,
		Role:  respUser.Role,
	}, nil
}

func (s *service) DeleteUser(ID uuid.UUID) error {
	_, err := s.repository.GetUserByID(ID)
	if err != nil {
		return errors.New("user does not exist")
	}

	return s.repository.DeleteUser(ID)
}
