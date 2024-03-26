package config

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"golang-authentication/internal/injector"
	"golang-authentication/internal/models"
	"gorm.io/gorm"
	"log"
)

type App struct {
	Fiber     *fiber.App
	database  *gorm.DB
	validator *validator.Validate
	viper     *viper.Viper
}

func errorHandlerConfig(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var status string
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	switch code {
	case 400:
		status = "Bad Request"
	case 401:
		status = "Unauthorized"
	case 403:
		status = "Forbidden"
	case 404:
		status = "Not Found"
	case 408:
		status = "Request Timeout"
	case 500:
		status = "Internal Server Error"

	}
	return ctx.Status(code).JSON(models.ErrorResponse{
		Status:  status,
		Message: e.Message,
	})
}

func NewApp(viper *viper.Viper, validator *validator.Validate, database *gorm.DB) *App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandlerConfig,
	})
	return &App{Fiber: app, viper: viper, validator: validator, database: database}

}

func (app *App) StartServer() {

	port := app.viper.GetInt("server.port")

	err := app.Fiber.Listen(fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatal("Error connecting to server %v", err)
	}
}

func (app *App) Setup() {
	signUpRoute := injector.InjectSignUpRoute(app.Fiber, app.database, app.validator)
	signUpRoute.Setup()

	authRoute := injector.InjectAuthRoute(app.Fiber, app.database, app.validator, app.viper)
	authRoute.Setup()
}
