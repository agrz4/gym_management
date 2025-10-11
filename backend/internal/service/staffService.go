package service

import (
	"errors"
	"gym_management/config"
	"gym_management/internal/models"
	"gym_management/internal/repository"

	"github.com/google/uuid"
)

type StaffService struct {
	repo repository.AuthRepository
}

func NewStaffService() *StaffService {
	return &StaffService{repo: repository.NewAuthRepository()}
}

// GetStaffs mengambil semua user dengan role 'staff'
func (s *StaffService) GetStaffs() ([]models.User, error) {
	var staffs []models.User
	// Gorm Find with Where clause for role
	if err := config.DB.Where("role IN (?)", []string{"staff", "admin"}).Find(&staffs).Error; err != nil {
		return nil, err
	}
	return staffs, nil
}

// CreateStaff
func (s *StaffService) CreateStaff(input models.RegisterInput, role string) (*models.User, error) {
	existingUser, _ := s.repo.FindByEmail(input.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("gagal hash password")
	}

	newStaff := models.User{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         role, // Bisa 'staff' atau 'admin'
		IsActive:     true,
	}

	if err := s.repo.Create(&newStaff); err != nil {
		return nil, err
	}
	return &newStaff, nil
}

// UpdateStaff
func (s *StaffService) UpdateStaff(id uuid.UUID, input models.RegisterInput) (*models.User, error) {
	staff, err := s.repo.FindByID(id)
	if err != nil || staff == nil || (staff.Role != "staff" && staff.Role != "admin") {
		return nil, errors.New("staff/admin tidak ditemukan")
	}

	staff.Name = input.Name
	staff.Email = input.Email // Hati-hati mengubah email, bisa melanggar unique constraint

	if err := s.repo.Update(staff); err != nil {
		return nil, errors.New("gagal memperbarui staff")
	}
	return staff, nil
}

// DeleteStaff
func (s *StaffService) DeleteStaff(id uuid.UUID) error {
	// Pastikan tidak menghapus diri sendiri atau admin utama (opsional)
	return config.DB.Where("id = ? AND role IN (?)", id, []string{"staff", "admin"}).Delete(&models.User{}).Error
}
