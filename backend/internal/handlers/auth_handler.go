package handlers

import (
	"gym_management/internal/models"
	"gym_management/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
