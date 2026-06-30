package services

import (
	"errors"
	"go-cloud/internal/dto"
	"mime/multipart"
	"os"
	"time"
)

type FileResponse struct {
	ID           uint      `json:"id"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"created_at"`
}

func (serv *Services) GetFiles(userID uint) ([]FileResponse, error) {

	files, err := serv.repo.GetFiles(userID)

	if err != nil {
		return []FileResponse{}, err
	}

	filesResponse := []FileResponse{}

	for _, file := range files {
		filesResponse = append(filesResponse, FileResponse{
			ID:           file.ID,
			OriginalName: file.OriginalName,
			MimeType:     file.MimeType,
			Size:         file.Size,
			CreatedAt:    file.CreatedAt,
		})
	}

	return filesResponse, nil

}

func (serv *Services) UploadFile(userID uint, fileHeader *multipart.FileHeader) error {

	const maxFileSize = 1 * 1024 * 1024 * 1024

	if fileHeader.Size <= 0 {
		return errors.New("file is empty")
	}

	if fileHeader.Size > maxFileSize {
		return errors.New("file is too large")
	}

	userDir, err := getUserDir(userID, true)

	if err != nil {
		return err
	}

	storageName, storagePath, err := saveFile(userDir, fileHeader)

	if err != nil {
		return err
	}

	input := dto.CreateFileInput{
		UserID:       userID,
		OriginalName: fileHeader.Filename,
		StorageName:  storageName,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		Size:         fileHeader.Size,
	}

	err = serv.repo.CreateFile(input)

	if err != nil {
		_ = os.Remove(storagePath)
		return err
	}

	return nil

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
