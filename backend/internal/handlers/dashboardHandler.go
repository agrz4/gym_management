package handlers

import (
	"gym_management/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var dashboardService = service.NewDashboardService()

// GetStatsHandler @route GET /api/dashboard/stats (Admin Only)
func GetStatsHandler(c *gin.Context) {
	stats, err := dashboardService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data statistik."})
		return
	}
	c.JSON(http.StatusOK, stats)
}
