package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yangjianhua/ginrest/model"
)

type GroupController struct {
	BaseController
}

func (this *GroupController) get(ctx *gin.Context) {
	sId := ctx.Param("id")
	id, _ := strconv.Atoi(sId)

	var g model.Group
	this.Context.DB.Where("id=?", id).First(&g)
	if g.ID > 0 {
		ctx.JSON(200, gin.H{"code": 0, "data": g})
	} else {
		ctx.JSON(200, gin.H{"code": -1, "msg": "Not Found"})
	}

}

func (this *GroupController) getPage(ctx *gin.Context) {
	var g []model.Group

	name := ctx.Query("name")
	description := ctx.Query("description")
	parent_id := ctx.Query("parent_id")

	p_no := ctx.DefaultQuery("page", "1")
	p_count := ctx.DefaultQuery("pagecount", "10")

	pageno, _ := strconv.Atoi(p_no)
	pagecount, _ := strconv.Atoi(p_count)

	tx := this.Context.DB.Where("1=1")
	if name != "" {
		tx = tx.Where("name LIKE ?", "%"+name+"%")
	}
	if description != "" {
		tx = tx.Where("description LIKE ?", "%"+description+"%")
	}
	if parent_id != "" {
		tx = tx.Where("parent_id = ?", parent_id)
	}

	p := Pagging(&PagingParam{
		DB:      tx,
		Page:    pageno,
		Limit:   pagecount,
		OrderBy: []string{"ID DESC"},
		ShowSQL: false,
	}, &g)

	ctx.JSON(200, p)
}

func (this *GroupController) create(ctx *gin.Context) {
	var g model.Group
	if ctx.ShouldBind(&g) == nil {
		if err := this.Context.DB.Create(&g); err.Error != nil {
			ctx.JSON(200, gin.H{"code": -1, "msg": err.Error.Error()})
			return
		} else {
			ctx.JSON(200, gin.H{"code": 0, "msg": "create OK", "data": g.ID})
			return
		}
	} else {
		ctx.JSON(200, gin.H{"code": -1, "msg": "Invalid Form Data."})
		return
	}
}

func (this *GroupController) update(ctx *gin.Context) {
	id := ctx.Param("id")

	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	parent_id := ctx.PostForm("parent_id")

	var n_pid = 0
	if parent_id != "" {
		n_pid, _ = strconv.Atoi(parent_id)
	}

	var g model.Group

	this.Context.DB.Where("id=?", id).First(&g)
	if g.ID > 0 {
		g.Name = name
		g.Description = description
		g.ParentId = (uint)(n_pid)

		if err := this.Context.DB.Save(&g); err.Error != nil {
			ctx.JSON(200, gin.H{"code": -1, "msg": "Error: " + err.Error.Error()})
			return
		} else {
			ctx.JSON(200, gin.H{"code": 0, "msg": "update successful"})
			return
		}
	} else {
		ctx.JSON(200, gin.H{"code": -1, "msg": "cannot find record"})
		return
	}
}

func (this *GroupController) delete(ctx *gin.Context) {
	id := ctx.Param("id")

	var g model.Group
	this.Context.DB.Where("id=?", id).First(&g)
	if g.ID > 0 {
		if err := this.Context.DB.Delete(&g); err.Error != nil {
			ctx.JSON(200, gin.H{"code": -1, "msg": "Error:" + err.Error.Error()})
		} else {
			ctx.JSON(200, gin.H{"code": 0, "msg": "删除成功"})
		}
	} else {
		ctx.JSON(200, gin.H{"code": -1, "msg": "Not Found"})
	}
}

func (this *GroupController) InitRouter() {
	this.Context = Ctx

	path := apiRootPath + "/group"
	this.AddToRouter(&Router{path: path, method: "GET", auth: true}, this.getPage)
	this.AddToRouter(&Router{path: path, method: "POST", auth: true}, this.create)
	this.AddToRouter(&Router{path: path + "/:id", method: "GET", auth: true}, this.get)
	this.AddToRouter(&Router{path: path + "/:id", method: "PUT", auth: true}, this.update)
	this.AddToRouter(&Router{path: path + "/:id", method: "DELETE", auth: true}, this.delete)
}
