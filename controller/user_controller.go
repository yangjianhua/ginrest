package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yangjianhua/ginrest/model"
)

type UserController struct {
	BaseController
}

func (this *UserController) get(ctx *gin.Context) {
	sId := ctx.Param("id")
	id, _ := strconv.Atoi(sId)

	var u model.User
	this.Context.DB.Where("id=?", id).Find(&u)

	if u == (model.User{}) {
		ctx.JSON(200, gin.H{"code": -1, "msg": "Not Found"})
	} else {
		ctx.JSON(200, gin.H{"code": 0, "data": u})
	}
}

func (this *UserController) getPage(ctx *gin.Context) {
	var u []model.User

	name := ctx.Query("name")
	email := ctx.Query("email")
	mobile := ctx.Query("mobile")
	tel := ctx.Query("tel")

	p_no := ctx.DefaultQuery("pageno", "1")
	p_count := ctx.DefaultQuery("pagecount", "10")

	pageno, _ := strconv.Atoi(p_no)
	pagecount, _ := strconv.Atoi(p_count)

	tx := this.Context.DB.Where("1=1")
	if name != "" {
		tx = tx.Where("name LIKE ?", "%"+name+"%")
	}
	if email != "" {
		tx = tx.Where("email LIKE ?", "%"+email+"%")
	}
	if mobile != "" {
		tx = tx.Where("mobile LIKE ?", "%"+mobile+"%")
	}
	if tel != "" {
		tx = tx.Where("tel LIKE ?", "%"+tel+"%")
	}

	p := Pagging(&PagingParam{
		DB:      tx,
		Page:    pageno,
		Limit:   pagecount,
		OrderBy: []string{"ID DESC"},
		ShowSQL: false,
	}, &u)

	ctx.JSON(200, p)
}

func (this *UserController) create(ctx *gin.Context) {
	var u model.User
	if ctx.ShouldBind(&u) == nil {
		// Encrypt the password before save.
		u.Password = GetBcrypt(u.Password)

		if err := this.Context.DB.Create(&u); err.Error != nil {
			ctx.JSON(200, gin.H{"code": -1, "msg": err.Error.Error()})
			return
		} else {
			ctx.JSON(200, gin.H{"code": 0, "msg": "create OK"})
			return
		}
	} else {
		ctx.JSON(200, gin.H{"code": -1, "msg": "Invalid Form Data."})
		return
	}
}

func (this *UserController) update(ctx *gin.Context) {
	id := ctx.Param("id") // From a URL param, not from post form
	// name := ctx.PostForm("name") // ignore name field, not update to table.

	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	mobile := ctx.PostForm("mobile")
	tel := ctx.PostForm("tel")

	var u model.User

	this.Context.DB.Where("id=?", id).First(&u)
	if u == (model.User{}) {
		ctx.JSON(200, gin.H{"code": -1, "msg": "cannot find record"})
		return
	}

	u.Email = email
	u.Mobile = mobile
	u.Tel = tel

	if password != "" {
		if len(password) < 4 {
			ctx.JSON(200, gin.H{"code": -1, "msg": "password must not less then 4 chars"})
			return
		}
		u.Password = GetBcrypt(password)
	}
	if err := this.Context.DB.Save(&u); err.Error != nil {
		ctx.JSON(200, gin.H{"code": -1, "msg": "Error: " + err.Error.Error()})
		return
	} else {
		ctx.JSON(200, gin.H{"code": 0, "msg": "update successful"})
		return
	}
}

func (this *UserController) delete(ctx *gin.Context) {
	sId := ctx.Param("id")
	id, _ := strconv.Atoi(sId)

	var u model.User
	this.Context.DB.Where("id=?", id).First(&u)
	if err := this.Context.DB.Delete(&u); err.Error != nil {
		ctx.JSON(200, gin.H{"code": -1, "msg": "Error:" + err.Error.Error()})
	} else {
		ctx.JSON(200, gin.H{"code": 0, "msg": "删除成功"})
	}
}

func (this *BaseController) changePwd(ctx *gin.Context) {
	id := ctx.Param("id")

	password := ctx.PostForm("password")
	if len(password) < 5 {
		ctx.JSON(200, gin.H{"code": -1, "msg": "No vaild password provided"})
		return
	}
	old_pwd := ctx.PostForm("oldpassword")

	if password == old_pwd {
		ctx.JSON(200, gin.H{"code": -1, "msg": "new password cannot be same with the old password"})
		return
	}

	var u model.User
	this.Context.DB.Where("id=?", id).First(&u)
	if u.ID > 0 {
		// Check old password if is correct
		if !MatchBcrypt(old_pwd, u.Password) {
			ctx.JSON(200, gin.H{"code": -1, "msg": "Old Password not match"})
			return
		}

		u.Password = GetBcrypt(password)
		if err := this.Context.DB.Save(&u); err.Error != nil {
			ctx.JSON(200, gin.H{"code": -1, "msg": "Error:" + err.Error.Error()})
		} else {
			ctx.JSON(200, gin.H{"code": 0, "msg": "修改密码成功"})
		}
	} else {
		ctx.JSON(200, gin.H{"code": -1, "msg": "User not found"})
	}
}

func (this *UserController) InitRouter() {
	this.Context = Ctx

	this.AddToRouter(&Router{path: "/api/user", method: "GET", auth: true}, this.getPage)
	this.AddToRouter(&Router{path: "/api/user", method: "POST", auth: true}, this.create)
	this.AddToRouter(&Router{path: "/api/user/:id", method: "GET", auth: true}, this.get)
	this.AddToRouter(&Router{path: "/api/user/:id", method: "PUT", auth: true}, this.update)
	this.AddToRouter(&Router{path: "/api/user/:id", method: "DELETE", auth: true}, this.delete)
	this.AddToRouter(&Router{path: "/api/user/changepwd/:id", method: "POST", auth: true}, this.changePwd)
}
