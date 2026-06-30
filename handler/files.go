package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetFiles(c *fiber.Ctx) error {

	id := c.Locals("userID").(uint)

	files, err := h.serv.GetFiles(id)

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
			"message": "file is required"})
	}

	err = h.serv.UploadFile(id, fileHeader)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "File uploaded successfully"})

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
