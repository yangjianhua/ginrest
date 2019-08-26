package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yangjianhua/ginrest/controller"
)

// 解决本地开发跨域问题
// 可以将 Access-Control-Allow-Origin 设置为CONFIG选项，以便保护服务器安全
func allowCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	context := controller.Context{}
	context.Init()
	defer context.Destory()

	if controller.CONFIG.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(allowCors())
	port := ":8080"
	if controller.CONFIG.ServerPort > 0 {
		port = ":" + strconv.Itoa(controller.CONFIG.ServerPort)
	}

	var c controller.BaseController

	c.InitializeRouter(r)
	log.Fatal(r.Run(port))
}
