package main

import "golang-authentication/internal/config"

func main() {
	viperConfig := config.NewViper("./../")
	validator := config.NewValidator()
	database := config.NewGorm(viperConfig)
	app := config.NewApp(viperConfig, validator, database)
	app.Setup()
	app.StartServer()

}
