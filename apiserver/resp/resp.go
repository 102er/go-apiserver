package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	RequestParameterError = " 9" //参数解析错误
)

func WriteParamError(ctx *gin.Context, err error) {
	WriteErrorApi(ctx, http.StatusBadRequest, 4, RequestParameterError, "Request parameter error", err)
}

func WriteObject(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func WriteError(ctx *gin.Context, httpStatus int, errorCode, message string, err error) {
	ctx.JSON(httpStatus, NewApiError(errorCode, message, 2, err))
}
func WriteErrorApi(ctx *gin.Context, httpStatus, showType int, errorCode, message string, err error) {
	ctx.JSON(httpStatus, NewApiError(errorCode, message, showType, err))
}

// WriteSuccess 适用于 post put delete 等操作类api
func WriteSuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ApiSuccess{
		Success: true,
	})
}

type ApiSuccess struct {
	Success bool   `json:"success"`
	TraceId string `json:"traceId,omitempty"`
	Host    string `json:"host,omitempty"`
}

type ApiError struct {
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	ErrorDetail  string `json:"errorDetail"`
	ShowType     int    `json:"showType,omitempty"`
}

func NewApiError(errorCode string, message string, showType int, err error) ApiError {
	return ApiError{}
}
