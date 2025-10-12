package main

import (
	"errors"
	"gym_management/config"
	"gym_management/internal/handlers"
	"gym_management/internal/models"
	"gym_management/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// InitialSetup runs database migrations/setup
func InitialSetup() {
	// Create PostgreSQL ENUM Type
	if err := config.DB.Exec("CREATE TYPE IF NOT EXISTS role_enum AS ENUM ('admin', 'staff', 'member')").Error; err != nil {
		// Log error jika ENUM tidak dapat dibuat (terutama jika sudah ada)
		log.Println("Role ENUM setup complete or error:", err)
	}

	// Auto Migrate Tables
	config.DB.AutoMigrate(&models.User{}, &models.GymPackage{}, &models.Attendance{})
	log.Println("Database tables auto-migrated successfully.")

	SeedData()
}

// SeedData inserts initial users and packages if they don't exist.
func SeedData() {
	// PENTING: Pastikan semua service yang dibutuhkan (NewStaffService, RegisterMemberService) tersedia.
	staffService := service.NewStaffService()

	// 1. Cek dan Buat Paket
	monthlyPkg := models.GymPackage{Name: "Bulanan", Price: 300000.00, DurationDays: 30, Benefits: "Akses 30 hari"}
	config.DB.Where(models.GymPackage{Name: "Bulanan"}).FirstOrCreate(&monthlyPkg)

	yearlyPkg := models.GymPackage{Name: "Tahunan", Price: 3000000.00, DurationDays: 365, Benefits: "Akses 1 tahun, gratis loker khusus"}
	config.DB.Where(models.GymPackage{Name: "Tahunan"}).FirstOrCreate(&yearlyPkg)

	// 2. Cek dan Buat Admin User
	var adminUser models.User
	if err := config.DB.Where("email = ?", "admin@gym.com").First(&adminUser).Error; errors.Is(err, gorm.ErrRecordNotFound) {

		adminInput := models.RegisterInput{
			Name:     "Super Admin",
			Email:    "admin@gym.com",
			Password: "securepassword",
		}

		// Buat Admin (menggunakan service yang mengizinkan role assignment)
		staffService.CreateStaff(adminInput, "admin")
		log.Println("Seed Data: Admin user created.")

		// Buat Staff
		staffInput := models.RegisterInput{
			Name:     "Staff Resepsionis",
			Email:    "staff@gym.com",
			Password: "securepassword",
		}
		staffService.CreateStaff(staffInput, "staff")
		log.Println("Seed Data: Staff user created.")

		// Buat Member
		memberInput := models.RegisterInput{
			Name:        "Member Contoh",
			Email:       "member@gym.com",
			Password:    "securepassword",
			PhoneNumber: "08123456789",
			PackageID:   &monthlyPkg.ID, // Gunakan ID paket yang sudah dibuat
		}
		// RegisterMemberService mengembalikan token, kita hanya ingin membuat user-nya di sini
		service.RegisterMemberService(memberInput)
		log.Println("Seed Data: Member user created.")
	}
}

func main() {
	godotenv.Load()
	config.ConnectDatabase()
	InitialSetup() // Jalankan Migrasi dan Seeding

	router := gin.Default()

	// Public Routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", handlers.LoginHandler)
		auth.POST("/register", handlers.RegisterMemberHandler)
		// Tambahkan Refresh Token dan Logout
		auth.POST("/refresh-token", handlers.RefreshTokenHandler)
		auth.POST("/logout", handlers.LogoutHandler)
	}

	// Protected Routes Group
	api := router.Group("/api")
	api.Use(handlers.AuthMiddleware())
	{
		// === ADMIN ONLY Routes ===
		admin := api.Group("/")
		admin.Use(handlers.RoleMiddleware("admin"))
		{
			// Member Delete
			admin.DELETE("/members/:id", handlers.DeleteMemberHandler)

			// Package CRUD (Full)
			admin.POST("/packages", handlers.CreatePackageHandler)
			admin.PUT("/packages/:id", handlers.UpdatePackageHandler)
			admin.DELETE("/packages/:id", handlers.DeletePackageHandler)

			// Staff CRUD (Full)
			admin.POST("/staff", handlers.CreateStaffHandler)
			admin.PUT("/staff/:id", handlers.UpdateStaffHandler)
			admin.DELETE("/staff/:id", handlers.DeleteStaffHandler)

			// Dashboard
			admin.GET("/dashboard/stats", handlers.GetStatsHandler)
		}

		// === ADMIN & STAFF Routes ===
		adminStaff := api.Group("/")
		adminStaff.Use(handlers.RoleMiddleware("admin", "staff"))
		{
			// Member Management (Read & Create/Update)
			adminStaff.GET("/members", handlers.GetMembersHandler)
			adminStaff.POST("/members", handlers.CreateMemberHandler)
			adminStaff.PUT("/members/:id", handlers.UpdateMemberHandler)

			// Attendance Operations
			adminStaff.POST("/attendance/checkin", handlers.CheckInHandler)
			adminStaff.POST("/attendance/checkout", handlers.CheckOutHandler)
			adminStaff.GET("/attendance/history", handlers.GetAllHistoryHandler)

			// Staff Read (Staff juga perlu melihat daftar staff)
			adminStaff.GET("/staff", handlers.GetStaffHandler)
		}

		// === PUBLIC (Authenticated) Routes ===
		// Paket bisa diakses oleh semua role yang sudah login
		api.GET("/packages", handlers.GetPackagesHandler)

		// === MEMBER ONLY Routes ===
		member := api.Group("/")
		member.Use(handlers.RoleMiddleware("member"))
		{
			// Member self-service
			member.GET("/attendance/my-history", handlers.GetMyHistoryHandler)
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
