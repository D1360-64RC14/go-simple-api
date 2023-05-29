package utils

import (
	"net/http"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/gin-gonic/gin"
)

func ErrorResponse(ctx *gin.Context, err *ErrorCode) {
	if err.Code() == http.StatusNotFound {
		ctx.Status(err.Code())
		return
	}

	ctx.JSON(err.Code(), dtos.NewErrorMessage(err))
}
