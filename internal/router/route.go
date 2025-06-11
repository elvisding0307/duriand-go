package router

import (
	account_handler "duriand/internal/handler/account"
	auth_handler "duriand/internal/handler/auth"
	ping_handler "duriand/internal/handler/ping"
	"duriand/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	// 需要JWT验证的API路由组
	v1Group := r.Group("/v1")
	{
		// 创建用户注册和登录的路由组
		authGroup := v1Group.Group("/auth")
		{
			authGroup.POST("/register", auth_handler.RegisterHandler)
			authGroup.POST("/login", auth_handler.LoginHandler)
		}

		// ping路由组
		pingGroup := v1Group.Group("/ping")
		pingGroup.Use(middleware.JWTAuth())
		{
			pingGroup.GET("", ping_handler.PingHandler)
		}

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
