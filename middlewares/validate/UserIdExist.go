package validate

import (
	"net/http"

	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/d1360-64rc14/simple-api/utils"
	"github.com/gin-gonic/gin"
)

func UserIdExist(userRepository interfaces.UserRepository) func(*gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.GetInt("id")

		exist, err := userRepository.UserExist(id)
		if err != nil {
			utils.ErrorAbortResponse(ctx, err)
			return
		}

		if !exist {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.Next()
	}
}
