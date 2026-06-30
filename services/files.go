package services

import (
	"errors"
	"go-cloud/internal/dto"
	"mime/multipart"
	"os"
)

func (serv *Services) GetFiles(userID uint, folderID *uint) ([]dto.FileResponse, error) {

	if folderID != nil {
		err := serv.repo.CheckFolderExists(userID, *folderID)

		if err != nil {
			return []dto.FileResponse{}, err
		}
	}

	files, err := serv.repo.GetFiles(userID, folderID)

	if err != nil {
		return []dto.FileResponse{}, err
	}

	filesResponse := []dto.FileResponse{}

	for _, file := range files {
		filesResponse = append(filesResponse, dto.FileResponse{
			ID:           file.ID,
			OriginalName: file.OriginalName,
			MimeType:     file.MimeType,
			Size:         file.Size,
			CreatedAt:    file.CreatedAt,
		})
	}

	return filesResponse, nil

}

func (serv *Services) UploadFile(userID uint, folderID *uint, fileHeader *multipart.FileHeader) (dto.FileResponse, error) {

	const maxFileSize = 1 * 1024 * 1024 * 1024

	emptyFileResponse := dto.FileResponse{}

	if fileHeader.Size <= 0 {
		return emptyFileResponse, errors.New("file is empty")
	}

	if fileHeader.Size > maxFileSize {
		return emptyFileResponse, errors.New("file is too large")
	}

	if folderID != nil {

		err := serv.repo.CheckFolderExists(userID, *folderID)

		if err != nil {
			return emptyFileResponse, err
		}

	}

	userDir, err := getUserDir(userID, true)

	if err != nil {
		return emptyFileResponse, err
	}

	storageName, storagePath, err := saveFile(userDir, fileHeader)

	if err != nil {
		return emptyFileResponse, err
	}

	input := dto.CreateFileInput{
		UserID:       userID,
		FolderID:     folderID,
		OriginalName: fileHeader.Filename,
		StorageName:  storageName,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		Size:         fileHeader.Size,
	}

	file, err := serv.repo.CreateFile(input)

	if err != nil {
		_ = os.Remove(storagePath)
		return emptyFileResponse, err
	}

	return file, nil

}

func (serv *Services) GetFileForDownload(userID uint, fileID uint) (string, string, error) {

	file, err := serv.repo.GetFileByIDAndUserID(fileID, userID)

	if err != nil {
		return "", "", err
	}

	userDir, err := getUserDir(userID, false)

	if err != nil {
		return "", "", err
	}

	storagePath := getStoragePath(userDir, file.StorageName)

	_, err = os.Stat(storagePath)

	if err != nil {
		return "", "", err
	}

	return storagePath, file.OriginalName, nil

}

func (serv *Services) DeleteFile(userID uint, fileID uint) error {

	file, err := serv.repo.GetFileByIDAndUserID(fileID, userID)

	if err != nil {
		return err
	}

	userDir, err := getUserDir(userID, false)

	if err != nil {
		return err
	}

	err = deleteFile(userDir, file.StorageName)

	if err != nil {
		return err
	}

	return serv.repo.DeleteFileByIDAndUserID(fileID, userID)

}

func (serv *Services) CreateFolder(userID uint, parentID *uint, name string) (dto.FolderResponse, error) {

	if parentID != nil {
		err := serv.repo.CheckFolderExists(userID, *parentID)

		if err != nil {
			return dto.FolderResponse{}, err
		}
	}

	input := dto.CreateFolderInput{
		UserID:   userID,
		ParentID: parentID,
		Name:     name,
	}

	return serv.repo.CreateFolder(input)

}

func (serv *Services) GetFolderContent(userID uint, folderID *uint) (dto.FolderContentResponse, error) {

	var currentFolder *dto.FolderResponse

	breadcrumbs := []dto.BreadcrumbItem{
		{ID: nil, Name: "My files"},
	}

	if folderID != nil {
		err := serv.repo.CheckFolderExists(userID, *folderID)

		if err != nil {
			return dto.FolderContentResponse{}, err
		}

		currentFolder, err = serv.repo.GetFolderByIDAndUserID(*folderID, userID)

		if err != nil {
			return dto.FolderContentResponse{}, err
		}

		breadcrumbsArray, err := serv.repo.GetBreadcrumbs(userID, *folderID)

		if err != nil {
			return dto.FolderContentResponse{}, err
		}

		breadcrumbs = append(breadcrumbs, breadcrumbsArray...)

	}

	files, err := serv.GetFiles(userID, folderID)

	if err != nil {
		return dto.FolderContentResponse{}, err
	}

	folders, err := serv.repo.GetFoldersByParentID(userID, folderID)

	if err != nil {
		return dto.FolderContentResponse{}, err
	}

	return dto.FolderContentResponse{
		Status:        "success",
		CurrentFolder: currentFolder,
		Breadcrumbs:   breadcrumbs,
		Files:         files,
		Folders:       folders,
	}, nil

}

func (serv *Services) UpdateFolderName(userID uint, folderID uint, newName string) (dto.FolderResponse, error) {

	err := serv.repo.CheckFolderExists(userID, folderID)

	if err != nil {
		return dto.FolderResponse{}, err
	}

	return serv.repo.UpdateFolderName(userID, folderID, newName)

}

func (serv *Services) DeleteFolder(userID uint, folderID uint) error {

	err := serv.repo.CheckFolderExists(userID, folderID)

	if err != nil {
		return err
	}

	folderIDs, err := serv.repo.GetFolderTreeIDs(userID, folderID)

	if err != nil {
		return err
	}

	files, err := serv.repo.GetFilesByFolderIDs(userID, folderIDs)

	if err != nil {
		return err
	}

	err = serv.repo.DeleteFolderTree(userID, folderIDs)

	if err != nil {
		return err
	}

	if len(files) > 0 {
		userDir, err := getUserDir(userID, false)

		if err != nil {
			return err
		}

		for _, file := range files {
			err = deleteFile(userDir, file.StorageName)
			if err != nil {
				return err
			}
		}
	}

	return nil

}
