package routers

import (
	"net/http"

	"github.com/d1360-64rc14/simple-api/middlewares"
	"github.com/gin-gonic/gin"
)

type Router struct {
	endpointPrefix string
	engine         *gin.Engine
	userController UserController
}

func NewRouter(endpointPrefix string, serverEngine *gin.Engine, userController UserController) *Router {
	router := &Router{
		endpointPrefix: endpointPrefix,
		engine:         serverEngine,
		userController: userController,
	}

	router.setup()

	return router
}

func (r Router) setup() {
	endpoint := r.engine.Group(r.endpointPrefix)

	endpoint.GET("/user", middlewares.ShouldHaveQuery("id"), r.userController.Get)
	endpoint.GET("/users", r.userController.GetAll)
	endpoint.POST("/user", r.userController.Create)
	endpoint.PUT("/user", middlewares.ShouldHaveQuery("id"), r.userController.Update)
	endpoint.DELETE("/user", middlewares.ShouldHaveQuery("id"), r.userController.Delete)

	endpoint.GET("/ping", r.ping)
}

func (r Router) ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong!")
}
