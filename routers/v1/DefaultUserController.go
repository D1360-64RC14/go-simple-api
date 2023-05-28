package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/d1360-64rc14/simple-api/config"
	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/gin-gonic/gin"
)

// DefaultUserController implements UserController
var _ interfaces.UserController = (*DefaultUserController)(nil)

type DefaultUserController struct {
	service interfaces.UserService
	repo    interfaces.UserRepository
}

func NewDefaultUserController(userService interfaces.UserService, userRepository interfaces.UserRepository) *DefaultUserController {
	return &DefaultUserController{
		service: userService,
		repo:    userRepository,
	}
}

func (c DefaultUserController) GetAll(ctx *gin.Context) {
	allUsers, err := c.service.SelectAllUsers()
	if err != nil {
		ctx.JSON(err.Code(), dtos.NewErrorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, allUsers)
}

func (c DefaultUserController) Get(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			dtos.NewErrorMessageByString(fmt.Sprintf("The user ID in the path should be an integer, not '%s'", idString)),
		)
		return
	}

	user, errC := c.service.SelectUserFromId(int(id))
	if errC != nil {
		ctx.Status(errC.Code())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c DefaultUserController) Create(ctx *gin.Context) {
	var newUser dtos.UserWithPassword

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorMessage(err))
		return
	}

	user, err := c.service.CreateUser(newUser)
	if err != nil {
		ctx.JSON(err.Code(), dtos.NewErrorMessage(err))
		return
	}

	newUserLocation := fmt.Sprintf("%s%s/%d", config.ApiUrl, ctx.Request.URL.Path, user.ID)
	ctx.Header("Location", newUserLocation)

	ctx.JSON(http.StatusCreated, user)
}

func (c DefaultUserController) Update(ctx *gin.Context) {
	ctx.Status(http.StatusNotImplemented)
}

func (c DefaultUserController) Delete(ctx *gin.Context) {
	ctx.Status(http.StatusNotImplemented)
}
