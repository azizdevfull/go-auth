package services

import (
	"go-auth/models"
	"go-auth/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, password string) error
	Authenticate(username, password string) (bool, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}
	return s.userRepo.CreateUser(user)
}

func (s *authService) Authenticate(username, password string) (bool, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
