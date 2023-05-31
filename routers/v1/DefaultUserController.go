package v1

import (
	"fmt"
	"net/http"

	"github.com/d1360-64rc14/simple-api/config"
	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/d1360-64rc14/simple-api/middlewares/validate"
	"github.com/d1360-64rc14/simple-api/utils"
	"github.com/gin-gonic/gin"
)

// DefaultUserController implements UserController
var _ interfaces.RouteController = (*DefaultUserController)(nil)

type DefaultUserController struct {
	service  interfaces.UserService
	repo     interfaces.UserRepository
	settings *config.Settings
}

func NewDefaultUserController(
	userService interfaces.UserService,
	userRepository interfaces.UserRepository,
	settings *config.Settings,
) interfaces.RouteController {
	return &DefaultUserController{
		service:  userService,
		repo:     userRepository,
		settings: settings,
	}
}

func (c DefaultUserController) AttachTo(group *gin.RouterGroup) {
	group.GET("/user/:id", validate.PathUserId, validate.UserIdExist(c.repo), c.get)
	group.GET("/users", c.getAll)
	group.POST("/user", c.create)
	group.PATCH("/user/:id", validate.PathUserId, validate.UserIdExist(c.repo), c.update)
	group.DELETE("/user/:id", validate.PathUserId, validate.UserIdExist(c.repo), c.delete)
}

func (c DefaultUserController) getAll(ctx *gin.Context) {
	allUsers, err := c.service.SelectAllUsers()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, allUsers)
}

func (c DefaultUserController) get(ctx *gin.Context) {
	id := ctx.GetInt("id")

	user, err := c.service.SelectUserFromId(id)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c DefaultUserController) create(ctx *gin.Context) {
	var newUser dtos.UserWithPassword

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorMessage(err))
		return
	}

	user, err := c.service.CreateUser(newUser)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	newUserLocation := fmt.Sprintf(
		"%s://%s%s/%d",
		c.settings.Api.Protocol,
		c.settings.Api.BaseUrl,
		ctx.Request.URL.Path,
		user.ID,
	)
	ctx.Header("Location", newUserLocation)

	ctx.JSON(http.StatusCreated, user)
}

func (c DefaultUserController) update(ctx *gin.Context) {
	id := ctx.GetInt("id")

	newUserData := new(dtos.UserUpdate)

	if err := ctx.ShouldBindJSON(&newUserData); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorMessage(err))
		return
	}

	err := c.service.UpdateUser(id, newUserData)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c DefaultUserController) delete(ctx *gin.Context) {
	id := ctx.GetInt("id")

	err := c.service.RemoveUser(id)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
