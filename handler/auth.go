package handler

import (
	"github.com/gofiber/fiber/v2"
)

type TrialAccessRequest struct {
	Email string `json:"email"`
}

func (h *Handler) SubmitTrialAccessRequest(c *fiber.Ctx) error {

	var request TrialAccessRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	successMsg, err := h.serv.SubmitTrialAccessRequest(request.Email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": successMsg})

}
