package register

import (
	"duriand/internal/dao"
	"duriand/internal/model"
	"errors"
	"time"
)

func RegisterService(username, password, corePassword string) error {
	tx := dao.DB_INSTANCE.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingUser model.User
	if err := tx.Where("username = ?", username).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return errors.New("user already exists")
	}

	user := model.User{
		Username:     username,
		Password:     password,
		CorePassword: corePassword,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	uid := user.Uid
	now := time.Now().Unix()
	// 更新timestamp
	timestamp := model.Timestamp{
		Uid:              uid,
		LatestUpdateTime: now,
		LatestDeleteTime: now,
	}
	if err := tx.Create(&timestamp).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
