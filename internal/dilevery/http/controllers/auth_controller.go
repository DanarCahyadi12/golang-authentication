package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang-authentication/internal/models"
	"golang-authentication/internal/usecase"
	"time"
)

type AuthController struct {
	AuthUseCase *usecase.AuthUseCase
}

func NewAuthController(authUseCase *usecase.AuthUseCase) *AuthController {
	return &AuthController{
		AuthUseCase: authUseCase,
	}
}

func (c *AuthController) SignIn(ctx *fiber.Ctx) error {
	body := new(models.SignInRequest)
	err := ctx.BodyParser(body)

	if err != nil {
		fmt.Println("Error parsing body ", err)
		return fiber.NewError(500, "Something wrong")
	}

	result, err := c.AuthUseCase.SignIn(ctx.Context(), body)
	if err != nil {
		fmt.Println("Error while sign in user: ", err)
		if e := err.(*models.ErrorResponse); e != nil {
			return fiber.NewError(e.Code, e.Message)
		}
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = result.RefreshToken
	cookie.Expires = time.Now().Add(3 * (24 * time.Hour))
	cookie.HTTPOnly = true

	ctx.Cookie(cookie)

	return ctx.Status(fiber.StatusOK).JSON(models.Response[*models.SignInResponse]{
		Message: "Sign in successfully",
		Data:    result,
	})

}

func (c *AuthController) GetToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token", "")
	result, err := c.AuthUseCase.GetToken(refreshToken)

	if err != nil {
		if e := err.(*models.ErrorResponse); e != nil {
			return fiber.NewError(e.Code, e.Message)
		} else {
			fmt.Println("Error", err)
			return fiber.NewError(500, "Something error")
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(models.Response[*models.GetTokenResponse]{Message: "Token successfully generated", Data: result})

}
