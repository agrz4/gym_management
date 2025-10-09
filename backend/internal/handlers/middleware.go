package handlers

import (
	"gym_management/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak. Token tidak valid."})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer")

		claims, err := service.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kadaluarsa."})
			return
		}

		// simpan claims ke context
		c.Set("userRole", claims.Role)
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

// RoleMiddleware membatasi akses berdasarkan role yang dibutuhkan
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Informasi role tidak ditemukan."})
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Tipe role tidak valid."})
			return
		}

		hasPermission := false
		for _, requiredRole := range requiredRoles {
			if roleStr == requiredRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses terlarang."})
			return
		}

		c.Next()
	}
}
