package middlewares

import (
	"net/http"
	"strings"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/gin-gonic/gin"
)

func ShouldHaveQuery(names ...string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		notExistingNames := make([]string, 0, len(names))

		for _, name := range names {
			if _, ok := ctx.GetQuery(name); !ok {
				notExistingNames = append(notExistingNames, name)
			}
		}

		if len(notExistingNames) == 0 {
			ctx.Next()
			return
		}

		errorMsg := "Query should have the following elements: " + strings.Join(notExistingNames, ", ")
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorMessageByString(errorMsg))
		ctx.Abort()
	}
}
