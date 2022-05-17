package utils

import (
	"github.com/gin-gonic/gin"
	"pipelineService/models/v1"
)

func BuildResponse(ctx *gin.Context, statusCode int, status string, err string, data interface{}) {
	response := models.Response{
		Status: status,
		Errors: err,
		Data:   data,
	}

	ctx.JSON(statusCode, response)
}

func BuildResponseAndAbort(ctx *gin.Context, statusCode int, status string, err string, data interface{}) {
	response := models.Response{
		Status: status,
		Errors: err,
		Data:   data,
	}

	ctx.AbortWithStatusJSON(statusCode, response)
}
