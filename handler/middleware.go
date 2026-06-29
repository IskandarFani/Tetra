package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TokenAuthDataRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) ReadRefreshToken(c *fiber.Ctx) error {

	var request TokenAuthDataRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body"})
	}

	if request.RefreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "missing authorization token",
		})
	}

	userID, userEmail, err := h.serv.UseRefreshToken(request.RefreshToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid refresh token",
		})
	}

	c.Locals("userID", userID)
	c.Locals("userEmail", userEmail)

	return c.Next()

}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)
	userEmail := c.Locals("userEmail").(string)

	tokens, err := h.serv.RefreshToken(userID, userEmail)

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

func (h *Handler) ReadToken(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "missing authorization header",
		})
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	accessToken = strings.TrimSpace(accessToken)

	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid access token"})
	}

	id, err := h.serv.ParseAccessToken(accessToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid access token"})
	}

	c.Locals("userID", id)

	return c.Next()

}
