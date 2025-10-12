package service

import (
	"gym_management/config"
	"gym_management/internal/models"

	"gorm.io/gorm"
)

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService() *DashboardService {
	return &DashboardService{db: config.DB}
}

type DashboardStats struct {
	TotalMembers            int64          `json:"totalMembers"`
	ActiveMembers           int64          `json:"activeMembers"`
	MembersByPackage        []PackageCount `json:"membersByPackage"`
	ProjectedMonthlyRevenue float64        `json:"projectedMonthlyRevenue"`
}

type PackageCount struct {
	PackageName string `json:"packageName"`
	MemberCount int64  `json:"memberCount"`
}

func (s *DashboardService) GetStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 1. Total & Aktif Member
	s.db.Model(&models.User{}).Where("role = ?", "member").Count(&stats.TotalMembers)
	s.db.Model(&models.User{}).Where("role = ? AND is_active = ?", "member", true).Count(&stats.ActiveMembers)

	// 2. Member per Paket (Menggunakan Raw SQL/Query Builder Gorm)
	rows, err := s.db.Model(&models.User{}).
		Select("gym_packages.name as package_name, count(users.id) as member_count").
		Joins("INNER JOIN gym_packages ON gym_packages.id = users.package_id").
		Where("users.role = ?", "member").
		Group("gym_packages.name").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pc PackageCount
		// Gorm scan menggunakan nama kolom dari Select
		if err := rows.Scan(&pc.PackageName, &pc.MemberCount); err != nil {
			return nil, err
		}
		stats.MembersByPackage = append(stats.MembersByPackage, pc)
	}

	// 3. Revenue Proyeksi (Simulasi: Total harga paket member aktif)
	// Catatan: Ini HANYA proyeksi. Revenue nyata harus dari tabel transaksi.
	var revenue float64
	s.db.Raw(`
		SELECT COALESCE(SUM(gp.price), 0) FROM users u 
		INNER JOIN gym_packages gp ON gp.id = u.package_id 
		WHERE u.role = 'member' AND u.is_active = TRUE
	`).Scan(&revenue)

	stats.ProjectedMonthlyRevenue = revenue

	return stats, nil
}
