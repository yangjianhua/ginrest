package controller

import (
	"log"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yangjianhua/ginrest/model"
)

type Router struct {
	method string
	regex  string
	path   string
	auth   bool
}

var Routers map[*Router]func(ctx *gin.Context)

type BaseController struct {
	// Session *model.Session
	Context *Context
	UserId  uint
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
	this.Context.DB.AutoMigrate(model.User{},
		model.Session{},
		model.Test{})

	ctx.JSON(200, gin.H{"code": 0, "msg": "Database Migrate OK."})
}

func (this *BaseController) InitRouter() {
	// Should be Called on every controller, to Init DB
	this.Context = Ctx

	// This is a Test Router just test for running.
	// this.AddToRouter(&Router{path: "welcome", method: "GET"}, welcome)

	// Init Json Web Token Login Object.
	this.jwtLogin()

	// This is a DB Migrate API
	this.AddToRouter(&Router{path: "/api/migrate", method: "POST"}, this.doMigrate)
	this.AddToRouter(&Router{path: "/api/login", method: "POST"}, authMiddleware.LoginHandler)
}

func (this *BaseController) getUserInfo(c *gin.Context) *model.User {
	claims := jwt.ExtractClaims(c)
	sId := claims[identityKey]
	var u model.User
	this.Context.DB.Where("id=?", sId).First(&u)

	return &u
}

// Init JWT Login Sample From https://github.com/appleboy/gin-jwt
func (this *BaseController) jwtLogin() {
	var err error
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
			id := claims["name"]
			this.Context.DB.Where("id=?", id).First(&u)
			return &u
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVars login
			if err := c.ShouldBind(&loginVars); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVars.Username
			password := loginVars.Password

			var u model.User
			this.Context.DB.Where("name=?", userID).First(&u)
			if MatchBcrypt(password, u.Password) {
				return &u, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*model.User); ok == true {
				return true
			}

			return false
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

func (this *BaseController) InitializeRouter(r *gin.Engine) {
	// All Controller's Router Register Here First.
	new(BaseController).InitRouter()
	new(TestController).InitRouter()
	new(UserController).InitRouter()

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		// claims := jwt.ExtractClaims(c)
		// log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "msg": "Page not found"})
	})

	apiPrefix := "/api"
	auth := r.Group(apiPrefix)
	auth.GET(apiPrefix+"/refresh_token", authMiddleware.RefreshHandler)

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
