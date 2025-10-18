package handlers

import (
	"gym_management/internal/models"
	"gym_management/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// api/auth/register
func RegisterMemberHandler(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "details": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := service.RegisterMemberService(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Registrasi berhasil.",
		"user":         gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role},
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// api/auth/login
func LoginHandler(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "details": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := service.LoginService(input.Email, input.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if strings.Contains(err.Error(), "aktif") {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login berhasil.",
		"user":         gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role},
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// @route POST /api/auth/refresh-token
// @access Public (with refresh token)
func RefreshTokenHandler(c *gin.Context) {
	var input models.RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token diperlukan."})
		return
	}

	accessToken, refreshToken, err := service.RefreshTokenService(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Token berhasil diperbarui.",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// @route POST /api/auth/logout
// @access Protected (any role)
func LogoutHandler(c *gin.Context) {
	// PENTING: Untuk logout, kita menggunakan AuthMiddleware untuk mendapatkan UserID
	// Namun, LogoutHandler berada di Public Routes Group, sehingga AuthMiddleware belum berjalan.
	// Kita harus memindahkan rute Logout ke Protected Group, atau menerima token dari body/header.
	// Pilihan yang lebih baik: Pindahkan ke Protected Group API.

	// Pindahkan rute '/logout' ke api.Group() di main.go

	// Asumsi rute ini berada di Protected Group (setelah AuthMiddleware):
	userID := c.MustGet("userID").(uuid.UUID)

	if err := service.LogoutService(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal logout."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil."})
}
