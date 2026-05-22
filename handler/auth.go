package handler

import (
	"github.com/gofiber/fiber/v2"
)

type NewUserStruct struct {
	Name string `json:"name"`
}

func (handlerStruct *Handler) AddNewUser(ctxStruct *fiber.Ctx) error {

	newUserNameRequestStruct := new(NewUserStruct)
	err := ctxStruct.BodyParser(newUserNameRequestStruct)

	if err != nil {
		return ctxStruct.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error()})
	}

	newUserName, errWithCreatingUser := handlerStruct.serv.AddNewUser(newUserNameRequestStruct.Name)

	if errWithCreatingUser != nil {
		return ctxStruct.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": errWithCreatingUser.Error()})
	}

	return ctxStruct.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": newUserName})

}
