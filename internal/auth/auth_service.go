package auth

import (
	"errors"
	"task-manager/internal/users"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(req RegisterUserRequest) (*users.User, error)
	LoginUser(req LoginUserRequest) (*users.User, error)
}

type authService struct {
	userRepo users.Repository
}

func NewAuthService(userRepo users.Repository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) RegisterUser(req RegisterUserRequest) (*users.User, error) {
	// 1. Check if user already exists
	existing, _ := s.userRepo.GetUserByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	// 2.Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Create user
	user, err := s.userRepo.CreateUser(&users.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user", // Default role is "user". Can be modified by an admin later.
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) LoginUser(req LoginUserRequest) (*users.User, error) {
	// 1. Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 3. Return user info (token generation will be handled in the handler)
	return user, nil
}
