package repository

import (
	"errors"
	"gym_management/config"
	"gym_management/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	FindUncheckedOutByUserID(userID uuid.UUID) (*models.Attendance, error)
	Create(attendance *models.Attendance) error
	Update(attendance *models.Attendance) error
	FindHistoryByUserID(userID uuid.UUID, limit int) ([]models.Attendance, error)
	FindAllHistory(filterUserID *uuid.UUID, dateFrom, dateTo *time.Time) ([]models.Attendance, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository() AttendanceRepository {
	return &attendanceRepository{db: config.DB}
}

// FindUncheckedOutByUserID: Mencari absensi hari ini yang belum CheckOut
func (r *attendanceRepository) FindUncheckedOutByUserID(userID uuid.UUID) (*models.Attendance, error) {
	var attendance models.Attendance

	todayStart := time.Now().Truncate(24 * time.Hour) // Awal hari ini

	err := r.db.Where("user_id = ? AND check_out_time IS NULL AND check_in_time >= ?", userID, todayStart).
		Order("check_in_time DESC").
		First(&attendance).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &attendance, nil
}

// Create implements AttendanceRepository.
func (r *attendanceRepository) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

// Update implements AttendanceRepository.
func (r *attendanceRepository) Update(attendance *models.Attendance) error {
	// Menggunakan Save() untuk Update berdasarkan Primary Key (ID)
	return r.db.Save(attendance).Error
}

// FindHistoryByUserID implements AttendanceRepository. (Untuk Member History)
func (r *attendanceRepository) FindHistoryByUserID(userID uuid.UUID, limit int) ([]models.Attendance, error) {
	var history []models.Attendance

	query := r.db.Where("user_id = ?", userID).
		Order("check_in_time DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

// FindAllHistory: Mengambil semua histori presensi dengan opsi filter.
func (r *attendanceRepository) FindAllHistory(filterUserID *uuid.UUID, dateFrom, dateTo *time.Time) ([]models.Attendance, error) {
	var history []models.Attendance

	// Preload User (Member) untuk mendapatkan nama/email
	query := r.db.Preload("User").Order("check_in_time DESC")

	// Filter berdasarkan User ID
	if filterUserID != nil && *filterUserID != uuid.Nil {
		query = query.Where("user_id = ?", *filterUserID)
	}

	// Filter berdasarkan Rentang Tanggal (CheckInTime)
	if dateFrom != nil {
		query = query.Where("check_in_time >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("check_in_time <= ?", *dateTo)
	}

	if err := query.Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}
