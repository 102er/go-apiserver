package controllers

import (
	v1 "github.com/102er/go-apiserver/api/v1"
	"github.com/102er/go-apiserver/apiserver/resp"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HelloWorldController struct {
	hwService v1.HelloWorldService
}

func RegisterApprovalChainHandler(router *gin.RouterGroup, hwService v1.HelloWorldService) {
	handler := &HelloWorldController{
		hwService,
	}
	rg := router.Group("hello-world")
	rg.GET("/", handler.Get)
}

// Get godoc
// @Summary 获取测试信息
// @Description 示例
// @Tags 示例
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} v1.ApprovalChain
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} resp.ApiError
// @Failure 500 {object} resp.ApiError
// @Router /api/v1/hello-world/{id} [get]
func (a *HelloWorldController) Get(ctx *gin.Context) {
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		resp.WriteParamError(ctx, err)
		return
	}
	ret, err := a.hwService.Get(ctx, id)
	if err != nil {
		resp.WriteError(ctx, http.StatusInternalServerError, "2", "hello world get failed", err)
		return
	}
	resp.WriteObject(ctx, ret)
}
