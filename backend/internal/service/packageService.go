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

func (s *PackageService) UpdatePackage(id uint, input models.UpdatePackageInput) (*models.GymPackage, error) {
	// 1. Cari Paket
	pkg, err := s.repo.FindByID(id)
	if err != nil || pkg == nil {
		return nil, errors.New("paket tidak ditemukan")
	}

	// 2. Update Fields (Cek apakah field di-supply, jika ya, update)
	if input.Name != "" {
		pkg.Name = input.Name
	}
	if input.Price > 0 {
		pkg.Price = input.Price
	}
	if input.DurationDays > 0 {
		pkg.DurationDays = input.DurationDays
	}
	// Benefits dapat berupa string kosong jika ingin menghapus benefit
	pkg.Benefits = input.Benefits

	// 3. Simpan ke Database
	if err := s.repo.Update(pkg); err != nil {
		return nil, errors.New("gagal memperbarui paket")
	}
	return pkg, nil
}

// DeletePackage handles business logic for deleting a gym package.
func (s *PackageService) DeletePackage(id uint) error {
	// 1. Cek keberadaan paket (opsional, repo.Delete akan menangani not found)
	// 2. Lakukan penghapusan
	if err := s.repo.Delete(id); err != nil {
		// Logika pengecekan Foreign Key Constraint Error dapat ditambahkan di sini
		// Jika Gorm gagal karena ada member yang menggunakan paket ini
		return errors.New("gagal menghapus paket. Mungkin masih ada member yang terikat dengan paket ini")
	}
	return nil
}
