package main

import (
	"gym_management/config"
	"gym_management/internal/handlers"
	"gym_management/internal/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// InitialSetup runs database migrations/setup
func InitialSetup() {
	// Create PostgreSQL ENUM Type
	if err := config.DB.Exec("CREATE TYPE role_enum AS ENUM ('admin', 'staff', 'member')").Error; err != nil {
		log.Println("Role ENUM already exists or cannot be created:", err)
	}

	// Auto Migrate Tables
	config.DB.AutoMigrate(&models.User{}, &models.GymPackage{}, &models.Attendance{})
	log.Println("Database tables auto-migrated successfully.")

	// Seed Data (TODO: Implement proper seeder logic here)
}

func main() {
	godotenv.Load()
	config.ConnectDatabase()
	InitialSetup()

	router := gin.Default()

	// public routes
	auth := router.Group("api/auth")
	{
		auth.POST("/login", handlers.LoginHandler)
		auth.POST("/register", handlers.RegisterMemberHandler)
	}

	// Protected Routes Group
	api := router.Group("/api")
	api.Use(handlers.AuthMiddleware())
	{
		// === ADMIN ONLY Routes ===
		admin := api.Group("/")
		admin.Use(handlers.RoleMiddleware("admin"))
		{
			// member delete (hanya admin)
			admin.DELETE("/members/:id", handlers.DeleteMemberHandler)

			// package CRUD
		}

		// === ADMIN & STAFF Routes ===
		adminStaff := api.Group("/")
		adminStaff.Use(handlers.RoleMiddleware("admin", "staff"))
		{
			// Example: Member Management
			// adminStaff.GET("/members", handlers.GetMembersHandler)
			// adminStaff.POST("/attendance/checkin", handlers.CheckInHandler)
		}

		// === MEMBER ONLY Routes ===
		member := api.Group("/")
		member.Use(handlers.RoleMiddleware("member"))
		{
			// Example: View Profile
			// member.GET("/member/profile", handlers.GetProfileHandler)
		}
	}

	// Server Start
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
