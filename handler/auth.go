package handler

import (
	"github.com/gofiber/fiber/v2"
)

type UserAuthDataRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(c *fiber.Ctx) error {

	var request UserAuthDataRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body"})
	}

	tokens, err := h.serv.Register(request.Email, request.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":        "success",
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken})

}

func (h *Handler) LogIn(c *fiber.Ctx) error {

	var request UserAuthDataRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body"})
	}

	tokens, err := h.serv.LogIn(request.Email, request.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":        "success",
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken})

}

func (h *Handler) Profile(c *fiber.Ctx) error {

	id := c.Locals("userID").(uint)

	userProfile, err := h.serv.GetUserProfile(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":         userProfile.ID,
		"email":      userProfile.Email,
		"created_at": userProfile.CreatedAt})

}
