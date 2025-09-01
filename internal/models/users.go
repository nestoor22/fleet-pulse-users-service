package models

import (
	"fleet-pulse-users-service/internal"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Email     string    `gorm:"not null;uniqueIndex"`
	Password  string
	internal.Metadata
}
