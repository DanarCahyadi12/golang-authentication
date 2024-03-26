package routes

import (
	"github.com/gofiber/fiber/v2"
	"golang-authentication/internal/dilevery/http/controllers"
)

type AuthRoute struct {
	App            *fiber.App
	AuthController *controllers.AuthController
}

func NewAuthRoute(app *fiber.App, controller *controllers.AuthController) *AuthRoute {
	return &AuthRoute{
		App:            app,
		AuthController: controller,
	}
}

func (r *AuthRoute) Setup() {
	r.App.Post("/auth", r.AuthController.SignIn)
	r.App.Get("/auth/token", r.AuthController.GetToken)
}
