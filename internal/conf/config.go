package conf

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// DURIAND 配置
	Host   string
	Port   string
	Secret string
	Expire int

	// MySQL 配置
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var DURIAND_CONFIG *Config = &Config{
	Host:       "0.0.0.0",
	Port:       "7224",
	Secret:     "123456",
	Expire:     86400,
	DBHost:     "127.0.0.1",
	DBPort:     "3306",
	DBUser:     "root",
	DBPassword: "123456",
	DBName:     "duriand",
}

func LoadConfig() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Fatalf("无法加载.env文件: %v", err)
	}
	cfg := DURIAND_CONFIG
	// 读取DURIAND配置
	if host := os.Getenv("DURIAND_HOST"); host != "" {
		cfg.Host = host
	}
	if port := os.Getenv("DURIAND_PORT"); port != "" {
		cfg.Port = port
	}
	if secret := os.Getenv("DURIAND_SECRET"); secret != "" {
		cfg.Secret = secret
	}
	if expireStr := os.Getenv("DURIAND_EXPIRE"); expireStr != "" {
		expire, err := strconv.Atoi(expireStr)
		if err != nil {
			log.Fatal("DURIAND_EXPIRE must be a valid integer")
		}
		cfg.Expire = expire
	}

	// 读取MySQL配置
	if dbHost := os.Getenv("MYSQL_HOST"); dbHost != "" {
		cfg.DBHost = dbHost
	}
	if dbPort := os.Getenv("MYSQL_PORT"); dbPort != "" {
		cfg.DBPort = dbPort
	}
	if dbUser := os.Getenv("MYSQL_USER"); dbUser != "" {
		cfg.DBUser = dbUser
	}
	if dbPassword := os.Getenv("MYSQL_PASSWORD"); dbPassword != "" {
		cfg.DBPassword = dbPassword
	}

	return
}
