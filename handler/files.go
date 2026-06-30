package handler

import (
	"go-cloud/internal/dto"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetFiles(c *fiber.Ctx) error {

	var folderID *uint

	id := c.Locals("userID").(uint)

	folderIDParam := c.Query("folder_id")

	if folderIDParam != "" {
		folderID64, err := strconv.ParseUint(folderIDParam, 10, 0)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid folder ID"})
		}
		folderIDUint := uint(folderID64)
		folderID = &folderIDUint
	}

	files, err := h.serv.GetFiles(id, folderID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(files)

}

func (h *Handler) UploadFile(c *fiber.Ctx) error {

	id := c.Locals("userID").(uint)
	fileHeader, err := c.FormFile("file")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "file is required",
		})
	}

	folderIDParam := c.FormValue("folder_id")

	var folderID *uint

	if folderIDParam != "" {
		folderID64, err := strconv.ParseUint(folderIDParam, 10, 0)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid folder ID",
			})
		}

		parsedFolderID := uint(folderID64)
		folderID = &parsedFolderID
	}

	file, err := h.serv.UploadFile(id, folderID, fileHeader)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "File uploaded successfully",
		"file":    file,
	})

}

func (h *Handler) DownloadFile(c *fiber.Ctx) error {

	id := c.Locals("userID").(uint)

	fileID, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid file ID"})
	}

	filePath, originalName, err := h.serv.GetFileForDownload(id, uint(fileID))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "File not found"})
	}

	return c.Download(filePath, originalName)

}

func (h *Handler) DeleteFile(c *fiber.Ctx) error {

	id := c.Locals("userID").(uint)

	fileID, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid file ID"})
	}

	err = h.serv.DeleteFile(id, uint(fileID))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "File deleted successfully"})

}

func (h *Handler) CreateFolder(c *fiber.Ctx) error {

	var request dto.CreateFolderRequest

	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body"})
	}

	id := c.Locals("userID").(uint)

	folderName := strings.TrimSpace(request.Name)

	if folderName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Folder name is required"})
	}

	folder, err := h.serv.CreateFolder(id, request.ParentID, folderName)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CreateFolderResponse{
		Status: "success",
		Folder: folder,
	})

}

func (h *Handler) GetFolderContent(c *fiber.Ctx) error {

	var folderID *uint

	id := c.Locals("userID").(uint)

	folderIDParam := c.Query("folder_id")

	if folderIDParam != "" {
		folderID64, err := strconv.ParseUint(folderIDParam, 10, 0)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid folder ID"})
		}
		folderIDUint := uint(folderID64)
		folderID = &folderIDUint
	}

	content, err := h.serv.GetFolderContent(id, folderID)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Folder not found"})
	}

	return c.Status(fiber.StatusOK).JSON(content)

}

func (h *Handler) UpdateFolder(c *fiber.Ctx) error {

	var request dto.UpdateFolderRequest

	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body"})
	}

	id := c.Locals("userID").(uint)

	folderID, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid folder ID"})
	}

	name := strings.TrimSpace(request.Name)

	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Folder name is required"})
	}

	folder, err := h.serv.UpdateFolderName(id, uint(folderID), name)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Folder updated successfully",
		"folder":  folder})

}

func (h *Handler) DeleteFolder(c *fiber.Ctx) error {

	id := c.Locals("userID").(uint)

	folderID, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid folder ID"})
	}

	err = h.serv.DeleteFolder(id, uint(folderID))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Folder deleted successfully"})
}
