package repository

import (
	"errors"
	"gym_management/config"
	"gym_management/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberRepository interface {
	FindAll(search string, isActive *bool) ([]models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(member *models.User) error
	Update(member *models.User) error
	Delete(id uuid.UUID) error
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository() MemberRepository {
	return &memberRepository{db: config.DB}
}

// FindAll implements MemberRepository.
func (r *memberRepository) FindAll(search string, isActive *bool) ([]models.User, error) {
	var members []models.User
	query := r.db.Where("role = ?", "member").Preload("Package")

	if search != "" {
		// ILIKE untuk case-insensitive di PostgreSQL
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// FindByID implements MemberRepository.
func (r *memberRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Package").First(&user, "id = ? AND role = ?", id, "member").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update implements MemberRepository.
func (r *memberRepository) Update(member *models.User) error {
	return r.db.Save(member).Error
}

// Delete implements MemberRepository.
func (r *memberRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ? AND role = ?", id, "member").Error
}

// Create dan FindByEmail sudah ada di auth_repo.go, tapi kita tetap perlu
// memastikan Create di sini untuk konsistensi, atau menggunakan AuthRepo untuk Create.
func (r *memberRepository) Create(member *models.User) error {
	return r.db.Create(member).Error
}

func (r *memberRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ? AND role = ?", email, "member").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
