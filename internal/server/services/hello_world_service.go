package services

import (
	v1 "github.com/102er/go-apiserver/api/v1"
	"github.com/gin-gonic/gin"
)

type HelloWorldService struct {
}

func newHelloWorldService() *HelloWorldService {
	return &HelloWorldService{}
}

func (h *HelloWorldService) Get(ctx *gin.Context, id int) (ret v1.HelloWorld, err error) {

	return
}
