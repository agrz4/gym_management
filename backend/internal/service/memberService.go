package service

import (
	"errors"
	"gym_management/internal/models"
	"gym_management/internal/repository"

	"github.com/google/uuid"
)

type MemberService struct {
	repo repository.MemberRepository
}

func NewMemberService() *MemberService {
	return &MemberService{repo: repository.NewMemberRepository()}
}

// createMember (untuk Admin/Staff)
func (s *MemberService) CreateMember(input models.RegisterInput) (*models.User, error) {
	if input.Email == "" || input.Password == "" {
		return nil, errors.New("email dan password wajib diisi")
	}

	existing, _ := s.repo.FindByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, _ := HashPassword(input.Password)

	member := models.User{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         "member",
		IsActive:     true,
		PhoneNumber:  input.PhoneNumber,
		Address:      input.Address,
		PackageID:    input.PackageID,
	}

	if err := s.repo.Create(&member); err != nil {
		return nil, errors.New("gagal menyimpan member ke database")
	}
	return &member, nil
}

// GetMembers (dengan filter/search)
func (s *MemberService) GetMembers(search string, isActive *bool) ([]models.User, error) {
	return s.repo.FindAll(search, isActive)
}

// UpdateMember
func (s *MemberService) UpdateMember(id uuid.UUID, input models.RegisterInput) (*models.User, error) {
	member, err := s.repo.FindByID(id)
	if err != nil || member == nil {
		return nil, errors.New("member tidak ditemukan")
	}

	// Update fields
	member.Name = input.Name
	member.PhoneNumber = input.PhoneNumber
	member.Address = input.Address
	member.PackageID = input.PackageID
	// isActive harus ditangani oleh input berbeda jika ada, atau tambahkan di struct input

	if err := s.repo.Update(member); err != nil {
		return nil, errors.New("gagal memperbarui member")
	}
	return member, nil
}

// DeleteMember
func (s *MemberService) DeleteMember(id uuid.UUID) error {
	return s.repo.Delete(id)
}
