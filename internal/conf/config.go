package conf

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// DURIAND 配置
	DuriandHost       string
	DuriandPort       string
	DuriandSecretKey  string
	DuriandExpireTime int

	// MySQL 配置
	MysqlHost     string
	MysqlPort     string
	MysqlUser     string
	MysqlPassword string
}

var DURIAND_CONFIG *Config = &Config{
	DuriandHost:       "0.0.0.0",
	DuriandPort:       "7224",
	DuriandSecretKey:  "123456",
	DuriandExpireTime: 86400,
	MysqlHost:         "127.0.0.1",
	MysqlPort:         "3306",
	MysqlUser:         "root",
	MysqlPassword:     "123456",
}

func LoadConfig() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Fatalf("无法加载.env文件: %v", err)
	}
	cfg := DURIAND_CONFIG
	// 读取DURIAND配置
	if host := os.Getenv("DURIAND_HOST"); host != "" {
		cfg.DuriandHost = host
	}
	if port := os.Getenv("DURIAND_PORT"); port != "" {
		cfg.DuriandPort = port
	}
	if secret := os.Getenv("DURIAND_SECRET_KEY"); secret != "" {
		cfg.DuriandSecretKey = secret
	}
	if expireStr := os.Getenv("DURIAND_EXPIRE_TIME"); expireStr != "" {
		expire, err := strconv.Atoi(expireStr)
		if err != nil {
			log.Fatal("DURIAND_EXPIRE must be a valid integer")
		}
		cfg.DuriandExpireTime = expire
	}

	// 读取MySQL配置
	if dbHost := os.Getenv("MYSQL_HOST"); dbHost != "" {
		cfg.MysqlHost = dbHost
	}
	if dbPort := os.Getenv("MYSQL_PORT"); dbPort != "" {
		cfg.MysqlPort = dbPort
	}
	if dbUser := os.Getenv("MYSQL_USER"); dbUser != "" {
		cfg.MysqlUser = dbUser
	}
	if dbPassword := os.Getenv("MYSQL_PASSWORD"); dbPassword != "" {
		cfg.MysqlPassword = dbPassword
	}

	return
}
