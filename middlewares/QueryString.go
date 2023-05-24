package middlewares

import (
	"net/http"
	"strings"

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

		errorMsg := "Query doesnt have elements " + strings.Join(notExistingNames, ", ")
		ctx.String(http.StatusBadRequest, errorMsg)
		ctx.Abort()
	}
}
