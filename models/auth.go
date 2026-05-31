package models

import "time"

type TrialRequest struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"type:varchar(255);not null;uniqueIndex"`
}

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash string `gorm:"type:varchar(255);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type RefreshToken struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"not null;index"`
	User   User `gorm:"foreignKey:UserID"`

	TokenHash string `gorm:"type:varchar(255);not null;uniqueIndex"`

	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
