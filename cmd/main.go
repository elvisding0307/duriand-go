package main

import (
	"duriand/internal/conf"
	"duriand/internal/dao"
	"duriand/internal/router"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Fatalf("无法加载.env文件: %v", err)
	}
	// 加载配置
	cfg := conf.LoadConfig()

	// 构建MySQL连接字符串
	mysqlURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)

	err := dao.InitDB(mysqlURL)
	if err != nil {
		panic("failed to connect database")
	}

	r := router.CreateRouter()

	// 使用配置中的端口
	addr := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("Starting Duriand on port %s\n", cfg.Port)
	err = r.Run(addr)
	if err != nil {
		panic("failed to start server")
	}
}
