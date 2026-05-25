package models

type TrialRequest struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"type:varchar(255);not null;uniqueIndex"`
}
