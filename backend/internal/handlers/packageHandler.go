package handlers

import (
	"gym_management/internal/models"
	"gym_management/internal/service"
	"net/http"
	"strconv"

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

// UpdatePackageHandler @route PUT /api/packages/:id (Admin Only)
func UpdatePackageHandler(c *gin.Context) {
	idParam := c.Param("id")
	packageID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID paket tidak valid."})
		return
	}

	var input models.UpdatePackageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid."})
		return
	}

	pkg, err := packageService.UpdatePackage(uint(packageID), input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Paket berhasil diperbarui.", "package": pkg})
}

// DeletePackageHandler @route DELETE /api/packages/:id (Admin Only)
func DeletePackageHandler(c *gin.Context) {
	idParam := c.Param("id")
	packageID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID paket tidak valid."})
		return
	}

	if err := packageService.DeletePackage(uint(packageID)); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // Menggunakan 409 Conflict untuk FK constraint
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Paket berhasil dihapus."})
}
