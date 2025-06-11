package router

import (
	account_handler "duriand/internal/handler/account"
	auth_handler "duriand/internal/handler/auth"
	login_handler "duriand/internal/handler/login"
	ping_handler "duriand/internal/handler/ping"
	register_handler "duriand/internal/handler/register"
	"duriand/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	apiGroup := r.Group("/api")
	// 需要JWT验证的API路由组
	v1Group := apiGroup.Group("/v1")
	{
		// 创建ping路由
		v1Group.GET("/ping", ping_handler.PingHandler)

		// 创建用户注册路由
		v1Group.POST("/register", register_handler.RegisterHandler)
		// 创建用户登录路由
		v1Group.POST("/login", login_handler.LoginHandler)

		// auth路由组
		authGroup := v1Group.Group("/auth")
		authGroup.Use(middleware.JWTAuth())
		{
			authGroup.GET("/verify", auth_handler.VerifyHandler)
		}

		// account路由组
		accountGroup := v1Group.Group("/account")
		accountGroup.Use(middleware.JWTAuth())
		{
			accountGroup.GET("", account_handler.QueryAccountHandler) // 修改为 accountGroup.GET("/query", api.QueryAccountHandler)
			accountGroup.POST("", account_handler.InsertAccountHandler)
			accountGroup.PUT("", account_handler.UpdateAccountHandler)
			accountGroup.DELETE("", account_handler.DeleteAccountHandler)
		}
	}

	return r
}
