package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`

	OriginalName string `gorm:"not null" json:"original_name"`
	StorageName  string `gorm:"not null;uniqueIndex" json:"-"`
	MimeType     string `gorm:"not null" json:"mime_type"`
	Size         int64  `gorm:"not null" json:"size"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
