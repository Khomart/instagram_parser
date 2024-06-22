package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type successResponseStruct struct {
	Message string `json:"message" example:"Message"`
	Data    any    `json:"data,omitempty"`
}

func validationError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, err.Error())
}

func invalidURL(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, "Only Instagram URLs are supported")
}

func successResponse(ctx *gin.Context, data any) {
	response := successResponseStruct{Message: "Success", Data: data}
	ctx.JSON(http.StatusOK, response)
}

func errorResponse(ctx *gin.Context, error error) {
	statusCode := http.StatusInternalServerError
	ctx.JSON(statusCode, error.Error())
}
