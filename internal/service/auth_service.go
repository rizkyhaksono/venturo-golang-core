package service

import (
	"context"
	"errors"
	"venturo-core/configs"
	"venturo-core/internal/model"
	"venturo-core/pkg/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db   *gorm.DB
	conf *configs.Config
}

// NewAuthService creates a new auth service.
func NewAuthService(db *gorm.DB, conf *configs.Config) *AuthService {
	return &AuthService{db: db, conf: conf}
}

// Register creates a new user.
func (s *AuthService) Register(ctx context.Context, name, email, password string) error {
	// Check if user already exists
	var existingUser model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&existingUser).Error; err == nil {
		return errors.New("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create new user
	newUser := model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Save user to the database
	if err := newUser.Save(s.db.WithContext(ctx)); err != nil {
		return err
	}

	return nil
}

// Login validates user credentials and returns access and refresh tokens.
func (s *AuthService) Login(ctx context.Context, email, password string) (map[string]string, error) {
	// Find user by email
	var user model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate access token (JWT)
	accessToken, err := utils.GenerateToken(user.ID, s.conf.JWTSecretKey)
	if err != nil {
		return nil, errors.New("could not generate access token")
	}

	// Generate refresh token
	refreshTokenString, err := model.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("could not generate refresh token")
	}

	// Hash the refresh token before storing
	hashedRefreshToken, err := model.HashToken(refreshTokenString)
	if err != nil {
		return nil, errors.New("could not hash refresh token")
	}

	// Remove existing refresh tokens for this user (single session per user)
	model.DeleteRefreshTokensByUserID(s.db.WithContext(ctx), user.ID)

	// Save new refresh token to database
	refreshToken := model.RefreshToken{
		UserID:      user.ID,
		HashedToken: hashedRefreshToken,
	}

	if err := refreshToken.Save(s.db.WithContext(ctx)); err != nil {
		return nil, errors.New("could not save refresh token")
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshTokenString,
		"token_type":    "Bearer",
	}, nil
}

// RefreshToken validates a refresh token and returns a new access token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (map[string]string, error) {
	if refreshTokenString == "" {
		return nil, errors.New("refresh token is required")
	}

	// Find all refresh tokens and check if any match
	var refreshTokens []model.RefreshToken
	if err := s.db.WithContext(ctx).Find(&refreshTokens).Error; err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var validRefreshToken *model.RefreshToken
	for _, rt := range refreshTokens {
		if model.ValidateToken(refreshTokenString, rt.HashedToken) {
			validRefreshToken = &rt
			break
		}
	}

	if validRefreshToken == nil {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(validRefreshToken.UserID, s.conf.JWTSecretKey)
	if err != nil {
		return nil, errors.New("could not generate access token")
	}

	return map[string]string{
		"access_token": accessToken,
		"token_type":   "Bearer",
	}, nil
}

// Logout invalidates the refresh token for a user.
func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	return model.DeleteRefreshTokensByUserID(s.db.WithContext(ctx), userID)
}
