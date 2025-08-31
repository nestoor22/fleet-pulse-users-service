package repositories

import (
	"fleet-pulse-users-service/internal"
	"fleet-pulse-users-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	*internal.BaseRepository[models.RefreshToken, uuid.UUID]
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	baseRepo := internal.NewBaseRepository[models.RefreshToken, uuid.UUID](db)
	return &RefreshTokenRepository{
		BaseRepository: baseRepo,
		db:             db,
	}
}

func (r *RefreshTokenRepository) GetByToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := r.db.Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) Update(tokenObj *models.RefreshToken) (*models.RefreshToken, error) {
	if err := r.db.Save(tokenObj).Error; err != nil {
		return nil, err
	}
	return tokenObj, nil
}

func (r *RefreshTokenRepository) Delete(tokenObj *models.RefreshToken) {
	r.db.Delete(tokenObj)
}

func (r *RefreshTokenRepository) DeletePreviousTokens(userID uuid.UUID) {
	r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{})
}
