package handlers

import (
	"gym_management/internal/models"
	"gym_management/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var memberService = service.NewMemberService()

// GET /api/members
func GetMembersHandler(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")

	var isActive *bool = nil
	if status != "" {
		active := status == "active"
		isActive = &active
	}

	members, err := memberService.GetMembers(search, isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data member."})
		return
	}
	c.JSON(http.StatusOK, members)
}

// POST /api/members
func CreateMemberHandler(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "details": err.Error()})
		return
	}

	member, err := memberService.CreateMember(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Member berhasil ditambahkan.", "member": member})
}

// PUT /api/members/:id
func UpdateMemberHandler(c *gin.Context) {
	idParam := c.Param("id")
	memberID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID member tidak valid."})
		return
	}

	var input models.RegisterInput // Menggunakan RegisterInput untuk update fields
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid.", "details": err.Error()})
		return
	}

	member, err := memberService.UpdateMember(memberID, input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member berhasil diperbarui.", "member": member})
}

// DELETE /api/members/:id
func DeleteMemberHandler(c *gin.Context) {
	idParam := c.Param("id")
	memberID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID member tidak valid."})
		return
	}

	if err := memberService.DeleteMember(memberID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus member."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member berhasil dihapus."})
}
