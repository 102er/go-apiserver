package v1

import (
	"github.com/102er/go-apiserver/internal/domain"
	"github.com/gin-gonic/gin"
)

type HelloWorldService interface {
	Get(ctx *gin.Context, id int) (ret HelloWorld, err error)
}

type HelloWorld struct {
	TypeMeta `json:",inline"`
	Data     struct {
		domain.HelloWorld
	} `json:"data"`
}

type HelloWorldList struct {
	TypeMeta `json:",inline"`
	Data     struct {
		Pages `json:",inline"`
		Data  []domain.HelloWorld `json:"data"`
	} `json:"data"`
}
