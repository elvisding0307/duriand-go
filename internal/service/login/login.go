package login

import (
	"duriand/internal/conf"
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

	tokenString, err := token.SignedString([]byte(conf.DURIAND_CONFIG.DuriandSecretKey))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}
