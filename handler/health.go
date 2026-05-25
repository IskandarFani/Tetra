package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CheckDB(c *fiber.Ctx) error {

	err := h.serv.CheckDB()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "database connection is healthy"})

}
