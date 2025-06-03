package main

import (
	"duriand/internal/conf"
	"duriand/internal/dao"
	"duriand/internal/router"
	"fmt"
)

func main() {
	// 加载配置
	conf.LoadConfig()
	cfg := conf.DURIAND_CONFIG
	err := dao.InitDB(cfg.MysqlHost, cfg.MysqlPort, cfg.MysqlUser, cfg.MysqlPassword)
	if err != nil {
		panic("failed to connect database")
	}

	r := router.CreateRouter()
	// 使用配置中的端口
	addr := fmt.Sprintf("%s:%s", cfg.DuriandHost, cfg.DuriandPort)
	fmt.Printf("Starting Duriand on port %s\n", cfg.DuriandPort)
	err = r.Run(addr)
	if err != nil {
		panic("failed to start server")
	}
}
