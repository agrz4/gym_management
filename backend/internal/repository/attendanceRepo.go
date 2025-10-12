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

// Create: Membuat record absensi baru
// func (r *attendanceRepository) Create(attendance *models.Attendance) error {
// 	return r.db.Create(attendance).Error
// }
