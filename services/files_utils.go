package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func getUserDir(userID uint, createDir bool) (string, error) {

	uploadDir := os.Getenv("UPLOAD_DIR")

	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	userDir := filepath.Join(uploadDir, fmt.Sprintf("user_%d", userID))

	if createDir {
		err := os.MkdirAll(userDir, 0755)

		if err != nil {
			return "", err
		}
	}

	return userDir, nil
}

func getStoragePath(userDir string, storageName string) string {

	storagePath := filepath.Join(userDir, storageName)

	return storagePath

}

func saveFile(userDir string, fileHeader *multipart.FileHeader) (string, string, error) {

	src, err := fileHeader.Open()

	if err != nil {
		return "", "", err
	}

	defer src.Close()

	ext := filepath.Ext(fileHeader.Filename)
	storageName := uuid.NewString() + ext

	storagePath := filepath.Join(userDir, storageName)

	dst, err := os.Create(storagePath)

	if err != nil {
		return "", "", err
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)

	if err != nil {
		_ = os.Remove(storagePath)
		return "", "", err
	}

	return storageName, storagePath, nil

}

func deleteFile(userDir string, storageName string) error {

	storagePath := getStoragePath(userDir, storageName)

	err := os.Remove(storagePath)

	if err != nil {

		if os.IsNotExist(err) {
			return nil
		}

		return err

	}

	return nil

}
