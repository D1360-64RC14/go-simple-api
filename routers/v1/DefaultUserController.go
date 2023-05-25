package v1

import (
	"net/http"
	"strconv"

	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/gin-gonic/gin"
)

// DefaultUserController implements UserController
var _ interfaces.UserController = (*DefaultUserController)(nil)

type DefaultUserController struct {
	service interfaces.UserService
}

func NewDefaultUserController(userService interfaces.UserService) *DefaultUserController {
	return &DefaultUserController{
		service: userService,
	}
}

func (c DefaultUserController) GetAll(ctx *gin.Context) {
	allUsers, err := c.service.SelectAllUsers()
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, allUsers)
}

func (c DefaultUserController) Get(ctx *gin.Context) {
	queryIdString, _ := ctx.GetQuery("id")
	queryId, _ := strconv.ParseInt(queryIdString, 10, 32)

	user, err := c.service.SelectUserFromId(int(queryId))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c DefaultUserController) Create(ctx *gin.Context) {
	username := ctx.PostForm("username")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	user, err := c.service.CreateUser(username, email, password)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c DefaultUserController) Update(ctx *gin.Context) {

}

func (c DefaultUserController) Delete(ctx *gin.Context) {

}
