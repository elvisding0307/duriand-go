package auth

import (
	"duriand/internal/dao"
	"duriand/internal/model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

func LoginService(username, password, corePassword string) (string, error) {
	var user model.User
	if err := dao.DB_INSTANCE.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	if user.Password != password || user.CorePassword != corePassword {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.Uid,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

func RegisterService(username, password, corePassword string) error {
	var existingUser model.User
	if err := dao.DB_INSTANCE.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return err
	}

	tx := dao.DB_INSTANCE.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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
	if err := tx.Model(&model.Timestamp{}).Create(&timestamp).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
