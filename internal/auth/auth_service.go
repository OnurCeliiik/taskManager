package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"task-manager/internal/users"
	"task-manager/utils/password"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(req RegisterUserRequest) (*users.User, error)
	LoginUser(req LoginUserRequest) (*users.User, error)
	ForgotPassword(email string) (string, error)
	ResetPassword(token, newPassword string) error
}

type authService struct {
	userRepo users.Repository
}

func NewAuthService(userRepo users.Repository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) RegisterUser(req RegisterUserRequest) (*users.User, error) {
	// 1. Validate password strength
	if err := password.Validate(req.Password); err != nil {
		return nil, err
	}

	// 2. Check if user already exists
	existing, _ := s.userRepo.GetUserByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	// 3.Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. Create user
	user, err := s.userRepo.CreateUser(&users.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      "user", // Default role is "user". Can be modified by an admin later.
		CreatedAt: time.Now(),
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

func (s *authService) ForgotPassword(email string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("if email exists, reset link will be sent")
	}

	resetToken := generateResetToken()
	expiryTime := time.Now().Add(1 * time.Hour)

	user.ResetToken = &resetToken
	user.ResetTokenExpiry = &expiryTime

	_, err = s.userRepo.UpdateUser(user)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *authService) ResetPassword(token, newPassword string) error {
	if err := password.Validate(newPassword); err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByResetToken(token)
	if err != nil {
		return errors.New("invalid reset token")
	}

	if user.ResetTokenExpiry == nil || time.Now().After(*user.ResetTokenExpiry) {
		return errors.New("reset token expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.ResetToken = nil
	user.ResetTokenExpiry = nil

	_, err = s.userRepo.UpdateUser(user)
	return err
}

func generateResetToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
