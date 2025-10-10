package repository

import (
	"gym_management/config"

	"gorm.io/gorm"
)

// PackageRepository Interface
type PackageRepository interface {
	// ... FindAll, FindByID, Create, Update, Delete
}

// packageRepository Implementation (Gorm)
type packageRepository struct {
	db *gorm.DB
}

// NewPackageRepository constructor
func NewPackageRepository() PackageRepository {
	return &packageRepository{db: config.DB}
}

// Implementasi CRUD di sini...
