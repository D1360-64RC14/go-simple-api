package routers

import (
	"net/http"

	"github.com/d1360-64rc14/simple-api/middlewares"
	"github.com/gin-gonic/gin"
)

// DefaultRouter implements Router
var _ Router = (*DefaultRouter)(nil)

type DefaultRouter struct {
	endpointPrefix string
	engine         *gin.Engine
	userController UserController
}

func NewDefaultRouter(endpointPrefix string, userController UserController) *DefaultRouter {
	router := &DefaultRouter{
		endpointPrefix: endpointPrefix,
		engine:         gin.Default(),
		userController: userController,
	}

	router.engine.Static("/api/docs", "routers/docs")
	router.engine.Use(middlewares.CORS)

	router.setup()

	return router
}

func (r DefaultRouter) setup() {
	endpoint := r.engine.Group(r.endpointPrefix)

	endpoint.GET("/user", middlewares.ShouldHaveQuery("id"), r.userController.Get)
	endpoint.GET("/users", r.userController.GetAll)
	endpoint.POST("/user", r.userController.Create)
	endpoint.PUT("/user", middlewares.ShouldHaveQuery("id"), r.userController.Update)
	endpoint.DELETE("/user", middlewares.ShouldHaveQuery("id"), r.userController.Delete)

	endpoint.GET("/ping", r.ping)
}

func (r DefaultRouter) ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong!")
}

func (r DefaultRouter) Engine() *gin.Engine {
	return r.engine
}

func (r DefaultRouter) EndpointPrefix() string {
	return r.endpointPrefix
}
