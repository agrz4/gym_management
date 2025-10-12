package models

import (
	"time"

	"github.com/google/uuid"
)

// --- DATABASE MODELS ---

// User Model (Digunakan untuk Admin, Staff, dan Member)
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:role_enum;default:'member';not null" json:"role"`

	// Member Specific
	PhoneNumber  string `gorm:"type:varchar(50)" json:"phoneNumber"`
	Address      string `gorm:"type:text" json:"address"`
	PackageID    *uint  `json:"packageId"`
	IsActive     bool   `gorm:"default:true" json:"isActive"`
	RefreshToken string `gorm:"type:text" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Package GymPackage `gorm:"foreignKey:PackageID" json:"package"`
}

type GymPackage struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	Name         string  `gorm:"type:varchar(100);unique;not null" json:"name"`
	Price        float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	DurationDays int     `gorm:"not null" json:"durationDays"`
	Benefits     string  `gorm:"type:text" json:"benefits"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Attendance struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"userId"`
	CheckInTime  time.Time  `gorm:"not null" json:"checkInTime"`
	CheckOutTime *time.Time `json:"checkOutTime"`

	User User `gorm:"foreignKey:UserID" json:"member"`
}

// --- INPUT STRUCTS ---

type RegisterInput struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	PackageID   *uint  `json:"packageId"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type CreatePackageInput struct {
	Name         string  `json:"name" binding:"required"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	DurationDays int     `json:"durationDays" binding:"required,gt=0"`
	Benefits     string  `json:"benefits"`
}

type UpdatePackageInput struct {
	Name         string  `json:"name"`
	Price        float64 `json:"price,omitempty"`
	DurationDays int     `json:"durationDays,omitempty"`
	Benefits     string  `json:"benefits"`
}
