package main

import (
	"flag"
	"github.com/102er/go-apiserver/apiserver/route"
	"github.com/102er/go-apiserver/pkg/logs"
	"github.com/gin-gonic/gin"
	"log"
)

// @title 102er接口文档
// @version 1.0
// @description 提供restful风格的api
// @termsOfService http://swagger.io/terms/

// @contact.name 102er
// @contact.url http://www.swagger.io/support
// @contact.email liaoxw@102er.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	configSource := flag.String("cs", "file", "config source")
	configPath := flag.String("config", "./config/config.yaml", "config file")
	flag.Parse()
	//load config
	if err := configer.LoadConfig(*configSource, *configPath); err != nil {
		log.Panic("web server init config failed,error:", err)
	}

	r := gin.New()
	r.Use(logs.GinLoggerMiddleware(), logs.GinRecoveryMiddleware(true))
	route.RegisterRoutes(r)
	if err := r.Run(":80"); err != nil {
		log.Panic("gin server run failed,error:", err)
	}
}
