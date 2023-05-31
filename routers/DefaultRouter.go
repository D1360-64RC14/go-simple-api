package routers

import (
	"net/http"

	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/d1360-64rc14/simple-api/middlewares"
	"github.com/gin-gonic/gin"
)

// DefaultRouter implements Router
var _ interfaces.Router = (*DefaultRouter)(nil)

type DefaultRouter struct {
	version          string
	endpointPrefix   string
	engine           *gin.Engine
	routeControllers []interfaces.RouteController
}

func NewDefaultV1Router(endpointPrefix string, routeControllers []interfaces.RouteController) interfaces.Router {
	version := "v1"

	router := &DefaultRouter{
		version:          version,
		endpointPrefix:   endpointPrefix + "/" + version,
		engine:           gin.Default(),
		routeControllers: routeControllers,
	}

	router.engine.Static(endpointPrefix+"/docs", "routers/docs")
	router.engine.Use(middlewares.CORS)

	router.setupRoutes()

	return router
}

func (r DefaultRouter) setupRoutes() {
	endpoint := r.engine.Group(r.endpointPrefix)
	endpoint.GET("/ping", r.ping)

	for _, group := range r.routeControllers {
		group.AttachTo(endpoint)
	}
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

func (r DefaultRouter) Version() string {
	return r.version
}
