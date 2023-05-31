package interfaces

import "github.com/gin-gonic/gin"

type RouteController interface {
	AttachTo(engine *gin.RouterGroup)
}
