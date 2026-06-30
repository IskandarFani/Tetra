package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID       uint  `gorm:"primaryKey"`
	UserID   uint  `gorm:"not null;index"`
	FolderID *uint `gorm:"index"`

	OriginalName string `gorm:"not null;size:255"`
	StorageName  string `gorm:"not null;size:255"`
	MimeType     string `gorm:"size:255"`
	Size         int64  `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Folder struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	ParentID  *uint  `gorm:"index"`
	Name      string `gorm:"not null;size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
