package controller

import (
	"log"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yangjianhua/ginrest/model"
	"gopkg.in/gormigrate.v1"
)

type Router struct {
	method string
	regex  string
	path   string
	auth   bool
	group  int
}

var Routers map[*Router]func(ctx *gin.Context)
var apiRootPath = "/api"

type BaseController struct {
	Context *Context
}

var identityKey = "id"
var authMiddleware *jwt.GinJWTMiddleware

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func init() {
	Routers = make(map[*Router]func(ctx *gin.Context))
}

func welcome(ctx *gin.Context) {
	// claims := jwt.ExtractClaims(ctx)

	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "welcome to ginrest",
		"data": "",
	})
}

func (this *BaseController) AddToRouter(r *Router, f gin.HandlerFunc) {
	Routers[r] = f
}

func (this *BaseController) doMigrate(ctx *gin.Context) {
	opt := &gormigrate.Options{
		TableName:      CONFIG.DbPrefix + "migrations",
		IDColumnName:   "id",
		IDColumnSize:   255,
		UseTransaction: true,
	}

	m := gormigrate.New(this.Context.DB, opt, []*gormigrate.Migration{
		{
			ID: "201810291200",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(model.User{}).Error; err != nil {
					return err
				}
				if err := tx.AutoMigrate(model.Group{}).Error; err != nil {
					return err
				}
				if err := tx.AutoMigrate(model.GroupUser{}).Error; err != nil {
					return err
				}

				// Init User.admin
				var u = model.User{
					Name:     "admin",
					Email:    "admin@admin.com",
					Password: GetBcrypt("admin123456"),
					IsAdmin:  true,
				}
				tx.Save(&u)

				// Add foreign keys
				tx.Model(model.GroupUser{}).AddForeignKey("user_id", CONFIG.DbPrefix+"users(id)", "CASCADE", "CASCADE")
				tx.Model(model.GroupUser{}).AddForeignKey("group_id", CONFIG.DbPrefix+"groups(id)", "CASCADE", "CASCADE")

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		ctx.JSON(200, gin.H{"code": -1, "msg": err.Error()})
	} else {
		ctx.JSON(200, gin.H{"code": 0, "msg": "DB Migreate OK."})
	}

}

func (this *BaseController) doCleanDb(ctx *gin.Context) {
	this.Context.DB.DropTableIfExists(&model.Test{}, &model.GroupUser{})

	this.Context.DB.DropTableIfExists(&model.User{}, &model.Group{})

	this.Context.DB.DropTableIfExists(CONFIG.DbPrefix + "migrations")

	ctx.JSON(200, gin.H{"code": 0, "msg": "database clean OK."})
}

func (this *BaseController) InitRouter() {
	// Should be Called on every controller, to Init DB
	this.Context = Ctx

	// This is a Test Router just test for running.
	this.AddToRouter(&Router{path: "/api/welcome", method: "GET"}, welcome)

	// Init Json Web Token Login Object.
	this.jwtLogin()

	// This is a DB Migrate API
	// The DB Migrate and Clean API should be disabled on product
	if CONFIG.Debug {
		this.AddToRouter(&Router{path: "/api/migrate", method: "POST"}, this.doMigrate)
		this.AddToRouter(&Router{path: "/api/cleandb", method: "POST", auth: true}, this.doCleanDb)
	}
	this.AddToRouter(&Router{path: "/api/login", method: "POST"}, authMiddleware.LoginHandler)
	this.AddToRouter(&Router{path: "/api/logout", method: "POST", auth: true}, this.logout)
	this.AddToRouter(&Router{path: "/api/userinfo", method: "GET", auth: true}, this.userInfo)
	this.AddToRouter(&Router{path: "/api/get_info", method: "GET", auth: true}, this.userInfo)
}

func (this *BaseController) getUserInfo(c *gin.Context) *model.User {
	claims := jwt.ExtractClaims(c)
	sId := claims[identityKey]
	var u model.User
	this.Context.DB.Where("id=?", sId).First(&u)

	return &u
}

func (this *BaseController) userInfo(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := claims[identityKey]
	var u model.User
	this.Context.DB.Where("id=?", id).First(&u)

	c.JSON(200, gin.H{"code": 0, "data": u})
}

// Init JWT Login Sample From https://github.com/appleboy/gin-jwt
func (this *BaseController) jwtLogin() {
	var err error
	var u_loc model.User // Save user for LoginResponse

	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour * 24 * 180,
		MaxRefresh:  time.Hour * 24 * 30,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"ID":        v.ID,
					"Name":      v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			var u model.User
			id := claims[identityKey]
			this.Context.DB.Where("id=?", id).First(&u)
			return &u
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVars login
			if err := c.ShouldBind(&loginVars); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVars.Username
			password := loginVars.Password

			var u model.User
			this.Context.DB.Where("name=?", username).First(&u)
			if MatchBcrypt(password, u.Password) {
				u_loc = u
				return &u, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if u, ok := data.(*model.User); ok && u.ID > 0 {
				if model.Uid != u.ID { // Set model.Uid Here, it's OK?
					model.Uid = u.ID
				}
				return true
			}

			return false
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"statusText": "OK",
				"token":      "Bearer " + token,
				"expire":     expire,
				"id":         u_loc.ID,
				"name":       u_loc.Name,
				"avator":     u_loc.Avator,
				"access":     "admin",
			})
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code": code,
				"msg":  message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
}

// Add logout code later.
func (this *BaseController) logout(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"code": 0, "msg": "Do nothing to logout."})
}

func (this *BaseController) InitializeRouter(r *gin.Engine) {
	// All Controller's Router Register Here First.
	new(BaseController).InitRouter()
	new(TestController).InitRouter()
	new(UserController).InitRouter()
	new(GroupController).InitRouter()
	new(GroupUserController).InitRouter()

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		// claims := jwt.ExtractClaims(c)
		// log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "msg": "Page not found"})
	})

	apiPrefix := apiRootPath
	auth := r.Group(apiPrefix)
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	// auth.Use(authMiddleware.MiddlewareFunc()).GET("/hello", welcome)

	// Loop to Init all routers.
	for k, v := range Routers {
		path := k.path
		if k.auth { // if auth, Only replace the very first "/api", the following "/api" will ignore
			startWith := strings.HasPrefix(path, apiPrefix)
			if startWith {
				path = strings.Replace(k.path, apiPrefix, "", 1) // So the controller can just write "/api/path", don't need to check auth
			}
		}

		// Add Router Right Here?
		// if k.group > 0 {
		// 	r.NoRoute(func(c *gin.Context) {
		// 		c.JSON(200, gin.H{"code": -1, "msg": "Not Allowed."})
		// 	})
		// 	continue
		// }

		switch k.method {
		case "GET":
			if k.auth {
				auth.Use(authMiddleware.MiddlewareFunc()).GET(path, v)
			} else {
				r.GET(path, v)
			}
			break
		case "POST":
			if k.auth {
				auth.Use(authMiddleware.MiddlewareFunc()).POST(path, v)
			} else {
				r.POST(path, v)
			}
			break
		case "PUT":
			if k.auth {
				auth.Use(authMiddleware.MiddlewareFunc()).PUT(path, v)
			} else {
				r.PUT(path, v)
			}
			break
		case "PATCH":
			if k.auth {
				auth.Use(authMiddleware.MiddlewareFunc()).PATCH(path, v)
			} else {
				r.PATCH(path, v)
			}
			break
		case "DELETE":
			if k.auth {
				auth.Use(authMiddleware.MiddlewareFunc()).DELETE(path, v)
			} else {
				r.DELETE(path, v)
			}
			break
		case "HEAD":
			if k.auth {
				auth.Use(authMiddleware.MiddlewareFunc()).DELETE(path, v)
			} else {
				r.HEAD(path, v)
			}
			break
		case "OPTIONS":
			if k.auth {
				auth.Use(authMiddleware.MiddlewareFunc()).OPTIONS(path, v)
			} else {
				r.OPTIONS(path, v)
			}
			break
		default:
		}
	}
}
