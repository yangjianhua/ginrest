package main

import (

	// "io"
	// "fmt"
	// "reflect"
	// "os"
	// "time"

	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yangjianhua/ginrest/controller"
)

// type dbContext struct {
// 	DB *gorm.DB
// }

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// func (this *dbContext) connectDB() {
// 	cnn := "root:root1234@tcp(127.0.0.1:3306)/ginrest?charset=utf8&parseTime=True&loc=Local"
// 	var err error = nil
// 	this.DB, err = gorm.Open("mysql", cnn)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// }

func main() {
	context := controller.Context{}
	context.Init()
	defer context.Destory()

	r := gin.Default()
	port := ":8080"
	if controller.CONFIG.ServerPort > 0 {
		port = ":" + strconv.Itoa(controller.CONFIG.ServerPort)
	}

	var c controller.BaseController

	c.InitializeRouter(r)
	log.Fatal(r.Run(port))
}
