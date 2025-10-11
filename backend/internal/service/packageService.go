package service

import (
	"errors"
	"gym_management/internal/models"
	"gym_management/internal/repository"
)

type PackageService struct {
	repo repository.PackageRepository
}

func NewPackageService() *PackageService {
	return &PackageService{repo: repository.NewPackageRepository()}
}

func (s *PackageService) GetPackages() ([]models.GymPackage, error) {
	return s.repo.FindAll()
}

func (s *PackageService) CreatePackage(input models.CreatePackageInput) (*models.GymPackage, error) {
	pkg := models.GymPackage{
		Name:         input.Name,
		Price:        input.Price,
		DurationDays: input.DurationDays,
		Benefits:     input.Benefits,
	}
	if err := s.repo.Create(&pkg); err != nil {
		return nil, errors.New("gagal membuat paket. Nama mungkin sudah ada.")
	}
	return &pkg, nil
}
