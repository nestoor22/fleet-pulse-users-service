package services

import (
	"fleet-pulse-users-service/internal/config"
	"fleet-pulse-users-service/internal/errors"
	"fleet-pulse-users-service/internal/models"
	"fleet-pulse-users-service/internal/repositories"
	"fleet-pulse-users-service/internal/schemas"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var inviteSecret = []byte(config.Get().Auth.InviteSecret)

type InviteClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	CompanyID uuid.UUID `json:"company_id"`
	jwt.RegisteredClaims
}

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	userRepo := repositories.NewUserRepository(db)
	return &UserService{repo: userRepo}
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
	var hashPassword string
	var err error
	if data.Password != "" {
		hashPassword, err = HashPassword(data.Password)
		if err != nil {
			return nil, err
		}
	}
	user, err := s.repo.Create(&models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  hashPassword,
	})
	if err != nil {
		return nil, err
	}
	err = s.SendInvite(user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) SendInvite(userID uuid.UUID) error {
	claims := InviteClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := token.SignedString(inviteSecret)
	fmt.Println(secret)
	return err
}

func (s *UserService) AcceptInvite(tokenString, password string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &InviteClaims{}, func(token *jwt.Token) (interface{}, error) {
		return inviteSecret, nil
	})
	if err != nil {
		return nil, errors.ErrInvalidInviteToken
	}

	claims, ok := token.Claims.(*InviteClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrInvalidInviteToken
	}

	user, err := s.GetUserById(claims.UserID)
	if user == nil || err != nil {
		return nil, errors.ErrUserNotFound
	}
	password, err = HashPassword(password)
	if err != nil {
		return nil, err
	}
	user, err = s.repo.Update(user, models.User{Password: password})
	if err != nil {
		return nil, err
	}
	return user, nil
}
