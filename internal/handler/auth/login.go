package auth

import (
	"duriand/internal/dao"
	"duriand/internal/handler"
	"duriand/internal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type LoginRequestSerializer struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CorePassword string `json:"core_password"`
}

func Login(c *gin.Context) {
	const (
		EMPTY_USERNAME_OR_PASSWORD int = iota + 1
		INVALID_CREDENTIALS
	)

	errorMap := map[int]string{
		EMPTY_USERNAME_OR_PASSWORD: "Username or password cannot be empty",
		INVALID_CREDENTIALS:        "Invalid username or password",
	}

	var req LoginRequestSerializer

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.CorePassword == "" {
		c.JSON(http.StatusOK, handler.NewErrorResponse(EMPTY_USERNAME_OR_PASSWORD, errorMap[EMPTY_USERNAME_OR_PASSWORD]))
		return
	}

	var user model.User
	if err := dao.DB_INSTANCE.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(INVALID_CREDENTIALS, errorMap[INVALID_CREDENTIALS]))
		return
	}
	// 验证密码和核心密码
	if user.Password != req.Password || user.CorePassword != req.CorePassword {
		c.JSON(http.StatusOK, handler.NewErrorResponse(INVALID_CREDENTIALS, errorMap[INVALID_CREDENTIALS]))
		return
	}

	// 生成 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.Uid,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // token 有效期24小时
	})

	// 使用密钥签名 token
	tokenString, err := token.SignedString([]byte("your-secret-key")) // 注意：实际使用时应该从配置文件读取密钥
	if err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(INVALID_CREDENTIALS, "Failed to generate token"))
		return
	}

	c.JSON(http.StatusOK, handler.NewSuccessResponse(map[string]string{
		"token": tokenString,
	}))
}
