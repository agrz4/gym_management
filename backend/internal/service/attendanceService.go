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

	// Cek apakah sudah Check-In
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
