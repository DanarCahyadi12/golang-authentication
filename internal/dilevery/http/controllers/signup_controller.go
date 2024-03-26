package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang-authentication/internal/models"
	"golang-authentication/internal/usecase"
)

type UserController struct {
	UserUseCase *usecase.SignUpUseCase
}

func NewUserController(useCase *usecase.SignUpUseCase) *UserController {
	return &UserController{
		UserUseCase: useCase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	body := new(models.SignUpRequest)
	err := ctx.BodyParser(body)
	if err != nil {
		fmt.Println("Error parsing body ", err)
		return fiber.NewError(500, "Something wrong")
	}

	result, err := c.UserUseCase.CreateUser(ctx.Context(), body)

	if err != nil {
		fmt.Println("Error while creating user: ", err)
		if e := err.(*models.ErrorResponse); e != nil {
			return fiber.NewError(e.Code, e.Message)
		}

		return fiber.NewError(500, "Something wrong with our server!")

	}
	return ctx.Status(fiber.StatusCreated).JSON(models.Response[*models.UserResponse]{Message: "Signup successfully", Data: result})

}
