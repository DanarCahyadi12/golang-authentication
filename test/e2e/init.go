package e2e

import (
	validator2 "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"golang-authentication/internal/config"
	"golang-authentication/internal/repository"
	"gorm.io/gorm"
)

var viperConfig *viper.Viper
var database *gorm.DB
var validator *validator2.Validate
var App *config.App
var UserRepository *repository.UserRepository

func init() {
	viperConfig = config.NewViper("./../../")
	validator = config.NewValidator()
	database = config.NewGorm(viperConfig)
	UserRepository = repository.NewUserRepository(database)
	App = config.NewApp(viperConfig, validator, database)
	App.Setup()

}
