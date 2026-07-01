package repository

import (
	"go-cloud/internal/dto"
	"go-cloud/models"

	"gorm.io/gorm"
)

func (repo *Repository) GetFiles(userID uint, folderID *uint) ([]models.File, error) {

	files := []models.File{}

	query := repo.db.Where("user_id = ?", userID)

	if folderID != nil {
		query = query.Where("folder_id = ?", *folderID)
	} else {
		query = query.Where("folder_id IS NULL")
	}

	err := query.Order("created_at DESC").Find(&files).Error

	return files, err

}

func (repo *Repository) CreateFile(params dto.CreateFileInput) (dto.FileResponse, error) {

	file := models.File{
		UserID:       params.UserID,
		FolderID:     params.FolderID,
		OriginalName: params.OriginalName,
		StorageName:  params.StorageName,
		MimeType:     params.MimeType,
		Size:         params.Size,
	}

	err := repo.db.Create(&file).Error

	if err != nil {
		return dto.FileResponse{}, err
	}

	return dto.FileResponse{
		ID:           file.ID,
		OriginalName: file.OriginalName,
		MimeType:     file.MimeType,
		Size:         file.Size,
		CreatedAt:    file.CreatedAt,
		FolderID:     file.FolderID,
	}, nil

}

func (repo *Repository) CreateFolder(params dto.CreateFolderInput) (dto.FolderResponse, error) {

	folder := models.Folder{
		UserID:   params.UserID,
		ParentID: params.ParentID,
		Name:     params.Name,
	}

	err := repo.db.Create(&folder).Error
	if err != nil {
		return dto.FolderResponse{}, err
	}

	return dto.FolderResponse{
		ID:        folder.ID,
		Name:      folder.Name,
		ParentID:  folder.ParentID,
		CreatedAt: folder.CreatedAt,
	}, nil

}

func (repo *Repository) GetFileByIDAndUserID(fileID uint, userID uint) (*models.File, error) {

	file := &models.File{}

	err := repo.db.Where("id = ? AND user_id = ?", fileID, userID).First(file).Error

	if err != nil {
		return nil, err
	}

	return file, nil

}

func (repo *Repository) CheckFolderExists(userID uint, folderID uint) error {

	folder := &models.Folder{}

	return repo.db.Where("user_id = ? AND id = ?", userID, folderID).First(folder).Error

}

func (repo *Repository) CheckFileExists(userID uint, fileID uint) error {

	file := &models.File{}

	return repo.db.Where("user_id = ? AND id = ?", userID, fileID).First(file).Error

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

func (repo *Repository) GetFolderByIDAndUserID(folderID uint, userID uint) (*dto.FolderResponse, error) {

	folder := &models.Folder{}

	err := repo.db.Where("id = ? AND user_id = ?", folderID, userID).First(folder).Error

	if err != nil {
		return nil, err
	}

	return &dto.FolderResponse{
		ID:        folder.ID,
		Name:      folder.Name,
		ParentID:  folder.ParentID,
		CreatedAt: folder.CreatedAt,
	}, nil

}

func (repo *Repository) GetFoldersByParentID(userID uint, parentID *uint) ([]dto.FolderResponse, error) {

	folders := []*models.Folder{}
	foldersArray := []dto.FolderResponse{}

	if parentID == nil {
		err := repo.db.Where("parent_id IS NULL AND user_id = ?", userID).Find(&folders).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := repo.db.Where("parent_id = ? AND user_id = ?", *parentID, userID).Find(&folders).Error
		if err != nil {
			return nil, err
		}
	}

	for _, row := range folders {
		foldersArray = append(foldersArray, dto.FolderResponse{
			ID:        row.ID,
			Name:      row.Name,
			ParentID:  row.ParentID,
			CreatedAt: row.CreatedAt,
		})
	}

	return foldersArray, nil
}

func (repo *Repository) GetBreadcrumbs(userID uint, folderID uint) ([]dto.BreadcrumbItem, error) {

	type breadcrumbRow struct {
		ID   uint
		Name string
	}

	breadcrumbs := []dto.BreadcrumbItem{}

	rows := []breadcrumbRow{}

	err := repo.db.Raw(`
		WITH RECURSIVE folder_path AS (
			SELECT
				id,
				name,
				parent_id,
				1 AS depth
			FROM folders
			WHERE id = ? AND user_id = ?

			UNION ALL

			SELECT
				f.id,
				f.name,
				f.parent_id,
				fp.depth + 1 AS depth
			FROM folders f
			JOIN folder_path fp ON f.id = fp.parent_id
			WHERE f.user_id = ?
		)
		SELECT id, name
		FROM folder_path
		ORDER BY depth DESC
	`, folderID, userID, userID).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		breadcrumbs = append(breadcrumbs, dto.BreadcrumbItem{
			ID:   &row.ID,
			Name: row.Name,
		})
	}

	return breadcrumbs, nil
}

func (repo *Repository) UpdateFolderName(userID uint, folderID uint, newName string) (dto.FolderResponse, error) {

	folder := &models.Folder{}

	err := repo.db.Where("id = ? AND user_id = ?", folderID, userID).First(folder).Error

	if err != nil {
		return dto.FolderResponse{}, err
	}

	folder.Name = newName

	err = repo.db.Save(folder).Error

	if err != nil {
		return dto.FolderResponse{}, err
	}

	return dto.FolderResponse{
		ID:        folder.ID,
		Name:      folder.Name,
		ParentID:  folder.ParentID,
		CreatedAt: folder.CreatedAt,
	}, nil

}

func (repo *Repository) GetFolderTreeIDs(userID uint, folderID uint) ([]uint, error) {

	var folderIDs []uint

	err := repo.db.Raw(`
		WITH RECURSIVE folder_tree AS (
			SELECT id
			FROM folders
			WHERE id = ? AND user_id = ?

			UNION ALL

			SELECT f.id
			FROM folders f
			JOIN folder_tree ft ON f.parent_id = ft.id
			WHERE f.user_id = ?
		)
		SELECT id FROM folder_tree
	`, folderID, userID, userID).Scan(&folderIDs).Error

	if err != nil {
		return nil, err
	}

	return folderIDs, nil

}

func (repo *Repository) GetFilesByFolderIDs(userID uint, folderIDs []uint) ([]models.File, error) {
	files := []models.File{}

	err := repo.db.
		Where("user_id = ? AND folder_id IN ?", userID, folderIDs).
		Find(&files).Error

	return files, err
}

func (repo *Repository) DeleteFolderTree(userID uint, folderIDs []uint) error {

	return repo.db.Transaction(func(tx *gorm.DB) error {

		err := tx.
			Where("user_id = ? AND folder_id IN ?", userID, folderIDs).
			Delete(&models.File{}).Error
		if err != nil {
			return err
		}

		err = tx.
			Where("user_id = ? AND id IN ?", userID, folderIDs).
			Delete(&models.Folder{}).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (repo *Repository) UpdateFileFolder(userID uint, fileID uint, folderID *uint) (dto.FileResponse, error) {

	file := &models.File{}

	err := repo.db.Where("user_id = ? AND id = ?", userID, fileID).First(file).Error

	if err != nil {
		return dto.FileResponse{}, err
	}

	file.FolderID = folderID

	err = repo.db.Save(file).Error

	if err != nil {
		return dto.FileResponse{}, err
	}

	return dto.FileResponse{
		ID:           file.ID,
		OriginalName: file.OriginalName,
		MimeType:     file.MimeType,
		Size:         file.Size,
		CreatedAt:    file.CreatedAt,
		FolderID:     file.FolderID,
	}, nil

}
