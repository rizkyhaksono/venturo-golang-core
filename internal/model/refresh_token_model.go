package model

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:char(36);not null" json:"user_id"`
	HashedToken string    `gorm:"size:255;not null;unique" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

// BeforeCreate is a GORM hook that runs before creating a new refresh token.
func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	rt.ID = uuid.New()
	return
}

// Save creates or updates a refresh token record.
func (rt *RefreshToken) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(rt).Error
}

// Delete removes a refresh token from the database.
func (rt *RefreshToken) Delete(db *gorm.DB) error {
	return db.WithContext(context.Background()).Delete(rt).Error
}

// GenerateRefreshToken creates a new random refresh token.
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken hashes a refresh token for secure storage.
func HashToken(token string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ValidateToken checks if a token matches the hashed version.
func ValidateToken(token, hashedToken string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
	return err == nil
}

// FindByUserID finds a refresh token by user ID.
func (rt *RefreshToken) FindByUserID(db *gorm.DB, userID uuid.UUID) error {
	return db.WithContext(context.Background()).Where("user_id = ?", userID).First(rt).Error
}

// DeleteByUserID removes all refresh tokens for a specific user.
func DeleteRefreshTokensByUserID(db *gorm.DB, userID uuid.UUID) error {
	return db.WithContext(context.Background()).Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}
