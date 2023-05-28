package validate

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/gin-gonic/gin"
)

func PathUserId(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.ParseInt(idString, 10, 32)
	if err == nil {
		ctx.Set("id", int(id))
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(
		http.StatusBadRequest,
		dtos.NewErrorMessageByString(fmt.Sprintf("The user ID in the path should be an integer, not '%s'", idString)),
	)
}
