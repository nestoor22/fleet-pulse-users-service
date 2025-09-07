package repositories

import (
	"fleet-pulse-users-service/internal"
	"fleet-pulse-users-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	*internal.BaseRepository[models.User, uuid.UUID]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	baseRepo := internal.NewBaseRepository[models.User, uuid.UUID](db)
	return &UserRepository{
		BaseRepository: baseRepo,
		db:             db,
	}
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetUsersByIDs(userIDs []uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
