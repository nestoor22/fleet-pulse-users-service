package internal

import "gorm.io/gorm"

type BaseRepository[T any, ID comparable] struct {
	db *gorm.DB
}

type BaseRepositoryInterface[T any, ID comparable] interface {
	GetById(id ID) (*T, error)
	Create(entity *T) (*T, error)
	Update(instance *T, inputData T) (*T, error)
	DeleteById(id ID) error
	DeleteObj(instance *T) error
}

var _ BaseRepositoryInterface[any, any] = &BaseRepository[any, any]{}

func NewBaseRepository[T any, ID comparable](db *gorm.DB) *BaseRepository[T, ID] {
	return &BaseRepository[T, ID]{db: db}
}

func (r *BaseRepository[T, ID]) GetById(id ID) (*T, error) {
	var entity T
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T, ID]) Create(entity *T) (*T, error) {
	if err := r.db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *BaseRepository[T, ID]) Update(instance *T, inputData T) (*T, error) {
	if err := r.db.Model(instance).Updates(inputData).Error; err != nil {
		return nil, err
	}
	return instance, nil
}

func (r *BaseRepository[T, ID]) DeleteById(id ID) error {
	var entity T
	return r.db.Delete(&entity, "id = ?", id).Error
}

func (r *BaseRepository[T, ID]) DeleteObj(instance *T) error {
	return r.db.Delete(instance).Error
}
