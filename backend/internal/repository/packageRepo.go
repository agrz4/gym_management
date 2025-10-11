package repository

import (
	"errors"
	"gym_management/config"
	"gym_management/internal/models"

	"gorm.io/gorm"
)

type PackageRepository interface {
	FindAll() ([]models.GymPackage, error)
	FindByID(id uint) (*models.GymPackage, error)
	Create(pkg *models.GymPackage) error
	Delete(id uint) error
}

type packageRepository struct {
	db *gorm.DB
}

// NewPackageRepository constructor
func NewPackageRepository() PackageRepository {
	return &packageRepository{db: config.DB}
}

func (r *packageRepository) FindAll() ([]models.GymPackage, error) {
	var pkgs []models.GymPackage
	if err := r.db.Order("price ASC").Find(&pkgs).Error; err != nil {
		return nil, err
	}
	return pkgs, nil
}

func (r *packageRepository) FindByID(id uint) (*models.GymPackage, error) {
	var pkg models.GymPackage
	if err := r.db.First(&pkg, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) Create(pkg *models.GymPackage) error {
	return r.db.Create(pkg).Error
}

func (r *packageRepository) Update(pkg *models.GymPackage) error {
	return r.db.Save(pkg).Error
}

func (r *packageRepository) Delete(id uint) error {
	// Gorm akan gagal jika ada foreign key constraint (member terkait), ini perilaku yang diinginkan
	return r.db.Delete(&models.GymPackage{}, id).Error
}
