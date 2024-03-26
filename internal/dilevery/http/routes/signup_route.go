package routes

import (
	"github.com/gofiber/fiber/v2"
	"golang-authentication/internal/dilevery/http/controllers"
)

type SignUpRoute struct {
	App *fiber.App
	*controllers.UserController
}

func NewUserRoute(app *fiber.App, controller *controllers.UserController) *SignUpRoute {
	return &SignUpRoute{
		App:            app,
		UserController: controller,
	}
}

func (r *SignUpRoute) Setup() {
	r.App.Post("/signup", r.UserController.Register)
}
