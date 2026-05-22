package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (handlerStruct *Handler) CheckDB(ctxStruct *fiber.Ctx) error {

	err := handlerStruct.serv.CheckDB()

	if err != nil {
		return ctxStruct.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	return ctxStruct.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "database connection is healthy"})

}
