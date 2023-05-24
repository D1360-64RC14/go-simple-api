package v1

import (
	"net/http"
	"strconv"

	"github.com/d1360-64rc14/simple-api/routers"
	"github.com/d1360-64rc14/simple-api/services"
	"github.com/gin-gonic/gin"
)

var _ routers.UserController = (*UserController)(nil)

type UserController struct {
	service services.UserService
}

func NewDefaultUserController(userService services.UserService) *UserController {
	return &UserController{
		service: userService,
	}
}

func (c UserController) GetAll(ctx *gin.Context) {
	allUsers, err := c.service.SelectAllUsers()
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, allUsers)
}

func (c UserController) Get(ctx *gin.Context) {
	queryIdString, _ := ctx.GetQuery("id")
	queryId, _ := strconv.ParseInt(queryIdString, 10, 32)

	user, err := c.service.SelectUserFromId(int(queryId))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c UserController) Create(ctx *gin.Context) {
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

func (c UserController) Update(ctx *gin.Context) {

}

func (c UserController) Delete(ctx *gin.Context) {

}
