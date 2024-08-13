package services

import (
	"go-auth/models"
	"go-auth/repositories"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, password string) error
	Authenticate(username, password string) (*models.User, error) // Updated return type
	CreateToken(id int, username string) (string, error)
	FindUserByUsername(username string) (*models.User, error) // New method
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *authService {
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

func (s *authService) Authenticate(username, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	// fmt.Println(user)
	return user, nil
}

var secretKey = []byte("secret-key")

func (s *authService) CreateToken(id int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       id,
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func (s *authService) FindUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetUserByUsername(username)
}
