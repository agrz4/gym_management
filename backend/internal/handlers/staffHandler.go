package handlers

import (
	"gym_management/internal/models"
	"gym_management/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var staffService = service.NewStaffService()

// GetStaffHandler @route GET /api/staff (Admin/Staff)
func GetStaffHandler(c *gin.Context) {
	staffs, err := staffService.GetStaffs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data staff."})
		return
	}
	c.JSON(http.StatusOK, staffs)
}

// CreateStaffHandler @route POST /api/staff (Admin Only)
func CreateStaffHandler(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid."})
		return
	}

	// Default role adalah 'staff'
	newStaff, err := staffService.CreateStaff(input, "staff")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Staff berhasil ditambahkan.", "staff": newStaff})
}

// UpdateStaffHandler @route PUT /api/staff/:id (Admin Only)
func UpdateStaffHandler(c *gin.Context) {
	idParam := c.Param("id")
	staffID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID staff tidak valid."})
		return
	}

	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid."})
		return
	}

	staff, err := staffService.UpdateStaff(staffID, input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff berhasil diperbarui.", "staff": staff})
}

// DeleteStaffHandler @route DELETE /api/staff/:id (Admin Only)
func DeleteStaffHandler(c *gin.Context) {
	idParam := c.Param("id")
	staffID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID staff tidak valid."})
		return
	}

	if err := staffService.DeleteStaff(staffID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus staff."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff berhasil dihapus."})
}
