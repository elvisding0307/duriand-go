package router

import (
	"duriand/internal/handler/api"
	"duriand/internal/handler/auth"
	"duriand/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	r := gin.Default()
	// 创建用户注册和登录的路由组
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", auth.RegisterHandler)
		authGroup.POST("/login", auth.LoginHandler)
	}

	// 需要JWT验证的API路由组
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.JWTAuth())
	{
		accountGroup := apiGroup.Group("/account")
		{
			accountGroup.GET("/query", api.QueryAccountHandler)
			accountGroup.POST("/insert", api.InsertAccountHandler)
			accountGroup.PUT("/update", api.UpdateAccountHandler)
			accountGroup.DELETE("/delete", api.DeleteAccountHandler)
		}
	}

	return r
}
