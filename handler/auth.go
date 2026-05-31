package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserAuthDataRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenAuthDataRequest struct {
	RefreshToken string `json:"refresh_token"`
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

func (h *Handler) RefreshToken(c *fiber.Ctx) error {

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

	tokens, err := h.serv.RefreshToken(request.RefreshToken)

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

	id, email, createdAt, err := h.serv.ParseAccessToken(accessToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid access token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":         id,
		"email":      email,
		"created_at": createdAt})

}
