package utils

import (
	"net/http"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/gin-gonic/gin"
)

func rule(err *ErrorCode, statusOnly func(int), statusJSON func(int, any)) {
	if err.Code() == http.StatusNotFound {
		statusOnly(err.Code())
		return
	}

	statusJSON(err.Code(), dtos.NewErrorMessage(err))
}

func ErrorResponse(ctx *gin.Context, err *ErrorCode) {
	rule(err, ctx.Status, ctx.JSON)
}

func ErrorAbortResponse(ctx *gin.Context, err *ErrorCode) {
	rule(err, ctx.AbortWithStatus, ctx.AbortWithStatusJSON)
}
