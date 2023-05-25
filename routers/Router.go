package routers

import "github.com/gin-gonic/gin"

type Router interface {
	EndpointPrefix() string
	Engine() *gin.Engine
}
