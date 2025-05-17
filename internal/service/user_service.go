package service

import (
	"errors"
	"habit-tracker-api/internal/auth"
	"habit-tracker-api/internal/domain"
	"habit-tracker-api/internal/repository"
	"habit-tracker-api/pkg/hash"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) Register(email, password string) error {
	hashed, err := hash.HashPassword(password)
	if err != nil {
		return err
	}
	user := &domain.User{
		Email:    email,
		Password: hashed,
	}
	return s.repo.Create(user)
}

var ErrInvalidCredentials = errors.New("invalid email or password")

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if !hash.CheckPasswordHash(password, user.Password) {
		return "", ErrInvalidCredentials
	}
	// Генерируем JWT
	token, err := auth.GenerateJWT(user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}
