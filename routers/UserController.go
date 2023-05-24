package routers

import "github.com/gin-gonic/gin"

type UserController interface {
	Get(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}
