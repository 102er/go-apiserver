package route

import (
	"github.com/102er/go-apiserver/apiserver/controllers"
	auth "github.com/102er/go-apiserver/internal/pkg/oauth"
	"github.com/102er/go-apiserver/internal/server/services"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func RegisterRoutes(r *gin.Engine) {
	//公共路由
	url := ginSwagger.URL("http://localhost:80/swagger/doc.json")
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	//不需要加权限控制的url

	//token校验 自定义中间件

	//版本路由

}

func RegisterRouterV1(r *gin.Engine) {
	//router := r.Group("/api/v1")
	routerJwt := r.Group("/api/v1", auth.JWTAuthMiddleware())
	controllers.RegisterApprovalChainHandler(routerJwt, services.BaseServiceInstance.HelloWorldService)
}
