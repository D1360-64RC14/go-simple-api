package v1

import (
	"net/http"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/gin-gonic/gin"
)

var _ interfaces.RouteController = (*DefaultAuthController)(nil)

type DefaultAuthController struct {
	userService interfaces.UserService
}

func NewDefaultAuthController(userService interfaces.UserService) interfaces.RouteController {
	return &DefaultAuthController{
		userService: userService,
	}
}

func (c DefaultAuthController) AttachTo(engine *gin.RouterGroup) {
	engine.GET("/auth", c.auth)
}

func (c DefaultAuthController) auth(ctx *gin.Context) {
	var authData dtos.AuthRequest

	if err := ctx.ShouldBindJSON(&authData); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorMessage(err))
		return
	}

}
