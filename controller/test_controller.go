package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/yangjianhua/ginrest/model"
)

type TestController struct {
	BaseController
}

func (this *TestController) create(ctx *gin.Context) {
	var t model.Test
	if ctx.ShouldBind(&t) == nil {
		if err := this.Context.DB.Create(&t); err.Error != nil {
			ctx.JSON(200, gin.H{"code": -1, "msg": "Error:" + err.Error.Error()})
			return
		} else {
			ctx.JSON(200, gin.H{"code": 0, "msg": "Create OK."})
			return
		}
	}
}

func (this *TestController) InitRouter() {
	this.Context = Ctx

	// This is a Test Router just test for running.
	// this.AddToRouter(&Router{path: "/api/test", method: "POST"}, this.create)
}
