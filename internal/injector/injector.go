package injector

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"golang-authentication/internal/dilevery/http/controllers"
	"golang-authentication/internal/dilevery/http/routes"
	"golang-authentication/internal/repository"
	"golang-authentication/internal/usecase"
	"gorm.io/gorm"
)

func InjectSignUpRoute(app *fiber.App, database *gorm.DB, validator *validator.Validate) *routes.SignUpRoute {
	userRepository := repository.NewUserRepository(database)
	userUseCase := usecase.NewSignupUseCase(userRepository, validator)
	userController := controllers.NewUserController(userUseCase)
	userRoute := routes.NewUserRoute(app, userController)
	return userRoute
}

func InjectAuthRoute(app *fiber.App, database *gorm.DB, validator *validator.Validate, viper *viper.Viper) *routes.AuthRoute {
	userRepository := repository.NewUserRepository(database)
	authUseCase := usecase.NewAuthUseCase(userRepository, validator, viper)
	authController := controllers.NewAuthController(authUseCase)
	authRoute := routes.NewAuthRoute(app, authController)

	return authRoute
}
