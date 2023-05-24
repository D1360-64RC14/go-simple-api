package middlewares

import "github.com/gin-contrib/cors"

var CORS = cors.New(
	cors.Config{
		AllowAllOrigins: true, // Development
		AllowMethods:    []string{"GET", "POST", "PATCH", "DELETE"},
	},
)
