package main

import (
	"duriand/internal/dao"
	"duriand/internal/router"
	"fmt"
)

func main() {
	// 替换为你的MySQL连接信息
	mysql_url := "root:123456@tcp(127.0.0.1:3307)/duriand?charset=utf8mb4&parseTime=True&loc=Local"

	err := dao.InitDB(mysql_url)
	if err != nil {
		panic("failed to connect database")
	}

	r := router.CreateRouter()

	fmt.Println("Starting Duriand on port 7224")
	err = r.Run(":7224")
	if err != nil {
		panic("failed to start server")
	}
}
