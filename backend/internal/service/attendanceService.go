package service

import (
	"errors"
	"gym_management/internal/models"
	"gym_management/internal/repository"
	"time"

	"github.com/google/uuid"
)

var (
	memberRepo     = repository.NewMemberRepository()
	attendanceRepo = repository.NewAttendanceRepository()
)

// CheckInMember: Hanya Staff yang bisa CheckIn
func CheckInMember(memberEmail string) (*models.Attendance, error) {
	member, err := memberRepo.FindByEmail(memberEmail)
	if err != nil || member == nil {
		return nil, errors.New("member tidak ditemukan")
	}
	if !member.IsActive {
		return nil, errors.New("member tidak aktif")
	}

	// Cek apakah sudah Check-In. Variabel err diabaikan (menggunakan _) karena
	// kita hanya peduli pada nilai existingAttendance (nil atau tidak nil)
	existingAttendance, _ := attendanceRepo.FindUncheckedOutByUserID(member.ID)
	if existingAttendance != nil {
		return nil, errors.New("member sudah Check-In dan belum Check-Out")
	}

	attendance := models.Attendance{
		UserID:      member.ID,
		CheckInTime: time.Now(),
	}

	if err := attendanceRepo.Create(&attendance); err != nil {
		return nil, errors.New("gagal menyimpan Check-In")
	}

	// Preload User/Member untuk response
	attendance.User = *member
	return &attendance, nil
}

// CheckOutMember: Hanya Staff yang bisa CheckOut
func CheckOutMember(memberEmail string) (*models.Attendance, error) {
	member, err := memberRepo.FindByEmail(memberEmail)
	if err != nil || member == nil {
		return nil, errors.New("member tidak ditemukan")
	}

	// Variabel err diabaikan karena kita hanya peduli pada nilai latestAttendance
	latestAttendance, _ := attendanceRepo.FindUncheckedOutByUserID(member.ID)
	if latestAttendance == nil {
		return nil, errors.New("member belum Check-In hari ini")
	}

	now := time.Now()
	latestAttendance.CheckOutTime = &now

	if err := attendanceRepo.Update(latestAttendance); err != nil {
		return nil, errors.New("gagal menyimpan Check-Out")
	}

	// Preload User/Member untuk response
	latestAttendance.User = *member
	return latestAttendance, nil
}

// GetMyHistory: Untuk member
func GetMyHistory(userID uuid.UUID) ([]models.Attendance, error) {
	return attendanceRepo.FindHistoryByUserID(userID, 50)
}

// GetAllHistory: Untuk Admin/Staff, mengelola filter dan memanggil repository.
func GetAllHistory(memberIDStr, dateFromStr, dateToStr string) ([]models.Attendance, error) {
	var filterUserID *uuid.UUID
	var dateFrom *time.Time
	var dateTo *time.Time

	// 1. Parsing Member ID (UUID)
	if memberIDStr != "" {
		if id, err := uuid.Parse(memberIDStr); err == nil {
			filterUserID = &id
		}
	}

	// 2. Parsing Tanggal Mulai
	if dateFromStr != "" {
		// Asumsi format ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)
		if t, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
			dateFrom = &t
		} else {
			// Jika parsing gagal, coba format tanggal saja (YYYY-MM-DD)
			if t, err := time.Parse("2006-01-02", dateFromStr); err == nil {
				dateFrom = &t
			}
		}
	}

	// 3. Parsing Tanggal Selesai
	if dateToStr != "" {
		if t, err := time.Parse(time.RFC3339, dateToStr); err == nil {
			dateTo = &t
		} else {
			if t, err := time.Parse("2006-01-02", dateToStr); err == nil {
				// Untuk dateTo, set waktu ke akhir hari untuk mencakup seluruh hari
				endOfDay := t.Add(24 * time.Hour).Add(-time.Second)
				dateTo = &endOfDay
			}
		}
	}

	return attendanceRepo.FindAllHistory(filterUserID, dateFrom, dateTo)
}
