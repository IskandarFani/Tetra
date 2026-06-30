package repository

import (
	"go-cloud/internal/dto"
	"go-cloud/models"

	"gorm.io/gorm"
)

func (repo *Repository) GetFiles(userID uint) ([]models.File, error) {

	files := []models.File{}

	err := repo.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&files).Error

	return files, err

}

func (repo *Repository) CreateFile(params dto.CreateFileInput) error {

	file := models.File{
		UserID:       params.UserID,
		OriginalName: params.OriginalName,
		StorageName:  params.StorageName,
		MimeType:     params.MimeType,
		Size:         params.Size,
	}

	return repo.db.Create(&file).Error

}

func (repo *Repository) GetFileByIDAndUserID(fileID uint, userID uint) (*models.File, error) {

	file := &models.File{}

	err := repo.db.Where("id = ? AND user_id = ?", fileID, userID).First(file).Error

	if err != nil {
		return nil, err
	}

	return file, nil

}

func (repo *Repository) DeleteFileByIDAndUserID(fileID uint, userID uint) error {

	result := repo.db.Where("id = ? AND user_id = ?", fileID, userID).Delete(&models.File{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}
