package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:role_enum;default:'member';not null" json:"role"` // PostgreSQL ENUM

	// Member Specific
	PhoneNumber  string `gorm:"type:varchar(50)" json:"phoneNumber"`
	Address      string `gorm:"type:text" json:"address"`
	PackageID    *uint  `json:"packageId"` // Foreign Key ke GymPackage
	IsActive     bool   `gorm:"default:true" json:"isActive"`
	RefreshToken string `gorm:"type:text" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
