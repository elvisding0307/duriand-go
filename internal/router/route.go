package router

import (
	"duriand/internal/controller"
	"duriand/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	r := gin.Default()
	// 创建用户注册和登录的路由组
	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
	}

	// 需要JWT验证的API路由组
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		api.POST("/insert_account", controller.InsertAccount)
		api.GET("/hello", controller.HelloWorld)
	}

	return r
}
