package controller

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/yangjianhua/ginrest/model"
)

type Context struct {
	DB   *gorm.DB
	User *model.User
}

var Ctx *Context

func init() {
	if Ctx == nil {
		Ctx = new(Context)
		Ctx.Init()
	}
}

func (this *Context) OpenDb() {
	var err error = nil
	this.DB, err = gorm.Open("mysql", CONFIG.DbUrl)
	if err != nil {
		str := fmt.Sprintf("Connect to Database Error, please check your config, %s", err.Error())
		panic(str)
	}

	// 赋值给model的全局变量，从而使model可以访问数据库
	// 不需要赋值，在model中可以使用 scope *gorm.Scope 获取当前的数据库访问链接
	// model.DB = this.DB
}

func (this *Context) CloseDb() {
	if this.DB != nil {
		this.DB.Close()
	}
}

func (this *Context) Init() {
	LoadConfig()
	this.OpenDb()
}

func (this *Context) Destory() {
	this.CloseDb()
}
