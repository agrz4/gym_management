package service

import (
	"errors"
	"gym_management/internal/models"
	"gym_management/internal/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

var authRepo = repository.NewAuthRepository()

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateTokens(user *models.User) (string, string, error) {
	// access token
	jwtExpiration := 24 // Default to 24 hours if not set
	if os.Getenv("JWT_EXPIRATION") != "" {
		if exp, err := time.ParseDuration(os.Getenv("JWT_EXPIRATION") + "h"); err == nil {
			jwtExpiration = int(exp.Hours())
		}
	}
	expirationTime := time.Now().Add(time.Duration(jwtExpiration) * time.Hour)
	claims := &AuthClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", "", err
	}

	// refresh token
	refreshExpiration := 168 // Default to 168 hours (7 days) if not set
	if os.Getenv("REFRESH_EXPIRATION") != "" {
		if exp, err := time.ParseDuration(os.Getenv("REFRESH_EXPIRATION") + "h"); err == nil {
			refreshExpiration = int(exp.Hours())
		}
	}
	refreshExpirationTime := time.Now().Add(time.Duration(refreshExpiration) * time.Hour)
	refreshClaims := &AuthClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	rt, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))

	return t, rt, err
}

// validate jwt token
func ValidateToken(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}

// registermembersevice handles member regis logic
func RegisterMemberService(input models.RegisterInput) (*models.User, string, string, error) {
	existingUser, _ := authRepo.FindByEmail(input.Email)
	if existingUser != nil {
		return nil, "", "", errors.New("email sudah terdaftar")
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		return nil, "", "", errors.New("gagal hash password")
	}

	newUser := models.User{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         "member", // Hardcoded role
		PhoneNumber:  input.PhoneNumber,
		Address:      input.Address,
		PackageID:    input.PackageID,
		IsActive:     true,
	}

	if err := authRepo.Create(&newUser); err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, _ := GenerateTokens(&newUser)
	newUser.RefreshToken = refreshToken
	authRepo.Update(&newUser)

	return &newUser, accessToken, refreshToken, nil
}

// login service handles user login
func LoginService(email, password string) (*models.User, string, string, error) {
	user, err := authRepo.FindByEmail(email)
	if err != nil {
		return nil, "", "", errors.New("kredensial tidak valid")
	}
	if user == nil || !CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", "", errors.New("kredensial tidak valid")
	}
	if !user.IsActive {
		return nil, "", "", errors.New("akun tidak aktif")
	}

	accessToken, refreshToken, _ := GenerateTokens(user)
	user.RefreshToken = refreshToken
	authRepo.Update(user)

	return user, accessToken, refreshToken, nil
}

// ValidateRefreshToken parses and validates the refresh token.
func ValidateRefreshToken(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}

// RefreshTokenService handles the logic for exchanging a refresh token for a new access token.
func RefreshTokenService(refreshTokenString string) (string, string, error) {
	claims, err := ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", errors.New("refresh token tidak valid atau kadaluarsa")
	}

	// 1. Cari User di DB
	user, err := authRepo.FindByID(claims.UserID)
	if err != nil {
		return "", "", errors.New("user tidak ditemukan")
	}

	// 2. Verifikasi apakah refresh token di DB cocok dengan yang dikirim
	if user.RefreshToken != refreshTokenString {
		// Ini mungkin tanda adanya serangan, sebaiknya semua token di-revoke
		user.RefreshToken = ""
		authRepo.Update(user)
		return "", "", errors.New("refresh token dicabut")
	}

	// 3. Generate token baru
	newAccessToken, newRefreshToken, err := GenerateTokens(user)
	if err != nil {
		return "", "", errors.New("gagal membuat token baru")
	}

	// 4. Update refresh token di DB
	user.RefreshToken = newRefreshToken
	authRepo.Update(user)

	return newAccessToken, newRefreshToken, nil
}

// LogoutService clears the refresh token in the database.
func LogoutService(userID uuid.UUID) error {
	user, err := authRepo.FindByID(userID)
	if err != nil {
		return nil // Jika user tidak ditemukan, anggap saja sudah logout
	}

	user.RefreshToken = ""
	return authRepo.Update(user)
}
