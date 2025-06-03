package conf

import (
	"os"
	"log"
	"strconv"
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

func LoadConfig() *Config {
	cfg := &Config{}

	// 读取DURIAND配置
	cfg.Host = os.Getenv("DURIAND_HOST")
	cfg.Port = os.Getenv("DURIAND_PORT")
	cfg.Secret = os.Getenv("DURIAND_SECRET")
	expireStr := os.Getenv("DURIAND_EXPIRE")
	if expireStr != "" {
		expire, err := strconv.Atoi(expireStr)
		if err != nil {
			log.Fatal("DURIAND_EXPIRE must be a valid integer")
		}
		cfg.Expire = expire
	}

	// 读取MySQL配置
	cfg.DBHost = os.Getenv("MYSQL_HOST")
	cfg.DBPort = os.Getenv("MYSQL_PORT")
	cfg.DBUser = os.Getenv("MYSQL_USER")
	cfg.DBPassword = os.Getenv("MYSQL_PASSWORD")
	cfg.DBName = os.Getenv("MYSQL_DATABASE")

	// 验证必填字段
	if cfg.Host == "" || cfg.Port == "" || cfg.Secret == "" {
		log.Fatal("DURIAND_HOST, DURIAND_PORT and DURIAND_SECRET are required")
	}
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		log.Fatal("MYSQL_HOST, MYSQL_USER and MYSQL_DATABASE are required")
	}

	return cfg
}