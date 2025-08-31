package services

import (
	"fleet-pulse-users-service/internal/errors"
	"fleet-pulse-users-service/internal/models"
	"fleet-pulse-users-service/internal/repositories"
	"fleet-pulse-users-service/internal/schemas"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserById(id uuid.UUID) (*models.User, error) {
	return s.repo.GetById(id)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) RegisterNewUser(data schemas.CreateUserRequest) (*models.User, error) {
	existingUser, _ := s.repo.GetUserByEmail(data.Email)
	if existingUser != nil {
		return nil, errors.ErrEmailAlreadyExists
	}
	hashPassword, err := HashPassword(data.Password)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(&models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  hashPassword,
	})
}
