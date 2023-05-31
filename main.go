package main

import (
	"log"

	"github.com/d1360-64rc14/simple-api/config"
	"github.com/d1360-64rc14/simple-api/database"
	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/d1360-64rc14/simple-api/repositories"
	"github.com/d1360-64rc14/simple-api/routers"
	v1 "github.com/d1360-64rc14/simple-api/routers/v1"
	"github.com/d1360-64rc14/simple-api/services"
)

func main() {
	settings, err := config.NewSettings("settings.yaml")
	fatalErr(err)

	database, err := database.NewMySQL(&settings.Database)
	fatalErr(err)

	userRepo, err := repositories.NewDefaultUserRepository(database)
	fatalErr(err)

	userService := services.NewDefaultUserService(userRepo)
	userController := v1.NewDefaultUserController(userService, userRepo, settings)

	controllers := []interfaces.RouteController{
		userController,
	}

	v1router := routers.NewDefaultV1Router("/api", controllers)
	v1router.Engine().Run(settings.Api.BaseUrl)
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
