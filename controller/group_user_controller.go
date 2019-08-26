package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/yangjianhua/ginrest/model"
)

type GroupUserController struct {
	BaseController
}

func (this *GroupUserController) get(ctx *gin.Context) {
	id := ctx.Param("id")

	var g model.GroupUser
	this.Context.DB.Where("id=?", id).First(&g)
	if g.ID > 0 {
		ctx.JSON(200, gin.H{"code": 0, "data": g})
	} else {
		ctx.JSON(200, gin.H{"code": -1, "msg": "Not Found"})
	}
}

func (this *GroupUserController) getPage(ctx *gin.Context) {

}

func (this *GroupUserController) create(ctx *gin.Context) {

}

func (this *GroupUserController) update(ctx *gin.Context) {

}

func (this *GroupUserController) delete(ctx *gin.Context) {

}

func (this *GroupUserController) InitRouter() {
	this.Context = Ctx

	path := apiRootPath + "/group_user"
	this.AddToRouter(&Router{path: path, method: "GET", auth: true}, this.getPage)
	this.AddToRouter(&Router{path: path, method: "POST", auth: true}, this.create)
	this.AddToRouter(&Router{path: path + "/:id", method: "GET", auth: true}, this.get)
	this.AddToRouter(&Router{path: path + "/:id", method: "PUT", auth: true}, this.update)
	this.AddToRouter(&Router{path: path + "/:id", method: "DELETE", auth: true}, this.delete)
}
