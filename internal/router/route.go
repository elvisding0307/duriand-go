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
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
	}

	// 需要JWT验证的API路由组
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.JWTAuth())
	{
		accountGroup := apiGroup.Group("/account")
		{
			accountGroup.GET("/query", api.QueryAccount)
			accountGroup.POST("/insert", api.InsertAccount)
			accountGroup.POST("/update", api.UpdateAccount)
		}
		apiGroup.GET("/hello", api.HelloWorld)
	}

	return r
}
