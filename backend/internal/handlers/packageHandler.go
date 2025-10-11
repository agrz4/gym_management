package handlers

import (
	"gym_management/internal/models"
	"gym_management/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var packageService = service.NewPackageService()

// GetPackagesHandler @route GET /api/packages (Public/Authenticated)
func GetPackagesHandler(c *gin.Context) {
	pkgs, err := packageService.GetPackages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil paket."})
		return
	}
	c.JSON(http.StatusOK, pkgs)
}

// CreatePackageHandler @route POST /api/packages (Admin Only)
func CreatePackageHandler(c *gin.Context) {
	var input models.CreatePackageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid."})
		return
	}

	pkg, err := packageService.CreatePackage(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pkg)
}
