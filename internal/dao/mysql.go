package dao

import (
	"duriand/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB_INSTANCE *gorm.DB

// InitDB 初始化数据库连接。
// 该函数应在整个应用程序生命周期中调用一次，以建立与数据库的连接。
// 参数:
//
//	mysql_connection_string (string): MySQL数据库的连接字符串，包含连接数据库所需的信息。
func InitDB(mysql_connection_string string) error {
	// 已经初始化了db
	if DB_INSTANCE != nil {
		return nil
	}
	db, err := Connect(mysql_connection_string)
	if err != nil {
		return err
	}
	DB_INSTANCE = db
	return nil
}

// Connect 尝试使用给定的连接字符串建立数据库连接，并自动迁移用户模型。
// 该函数接受一个连接字符串作为参数，返回一个*gorm.DB实例和一个错误对象。
// 如果连接成功，错误对象将为nil；如果连接失败，将返回nil和错误对象。
func Connect(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移所有模型
	return db, db.AutoMigrate(&model.User{}, &model.Account{}, &model.Timestamp{})
}
