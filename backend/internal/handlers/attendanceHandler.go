package handlers

import (
	"gym_management/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttendanceInput struct {
	MemberEmail string `json:"memberEmail" binding:"required,email"`
}

// CheckInHandler @route POST /api/attendance/checkin (Staff Only)
func CheckInHandler(c *gin.Context) {
	var input AttendanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email member diperlukan."})
		return
	}

	attendance, err := service.CheckInMember(input.MemberEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    attendance.User.Name + " berhasil Check-In!",
		"attendance": attendance,
	})
}

// CheckOutHandler @route POST /api/attendance/checkout (Staff Only)
func CheckOutHandler(c *gin.Context) {
	var input AttendanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email member diperlukan."})
		return
	}

	attendance, err := service.CheckOutMember(input.MemberEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    attendance.User.Name + " berhasil Check-Out!",
		"attendance": attendance,
	})
}

// GetMyHistoryHandler @route GET /api/attendance/my-history (Member Only)
func GetMyHistoryHandler(c *gin.Context) {
	// userID diambil dari JWT claims
	userID := c.MustGet("userID").(uuid.UUID)

	history, err := service.GetMyHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil histori."})
		return
	}
	c.JSON(http.StatusOK, history)
}

// GetAllHistoryHandler @route GET /api/attendance/history (Admin/Staff Only)
func GetAllHistoryHandler(c *gin.Context) {
	// Mengambil query parameters untuk filter
	memberID := c.Query("member_id")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Panggil service untuk mendapatkan histori dengan filter
	history, err := service.GetAllHistory(memberID, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil semua histori presensi."})
		return
	}

	c.JSON(http.StatusOK, history)
}
