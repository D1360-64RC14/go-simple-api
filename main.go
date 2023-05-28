package main

import (
	"database/sql"
	"log"

	_ "github.com/proullon/ramsql/driver" // Needed to ramsql work

	"github.com/d1360-64rc14/simple-api/repositories"
	"github.com/d1360-64rc14/simple-api/routers"
	v1 "github.com/d1360-64rc14/simple-api/routers/v1"
	"github.com/d1360-64rc14/simple-api/services"
)

func main() {
	db, err := sql.Open("ramsql", "database")
	fatalErr(err)

	userRepo, err := repositories.NewDefaultUserRepository(db)
	fatalErr(err)

	userService := services.NewDefaultUserService(userRepo)
	userController := v1.NewDefaultUserController(userService, userRepo)

	router := routers.NewDefaultRouter("/api/v1", userController)
	router.Engine().Run("localhost:1360")
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
