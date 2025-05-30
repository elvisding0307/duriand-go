package auth

import (
	"duriand/internal/dao"
	"duriand/internal/handler"
	"duriand/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequestSerializer struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CorePassword string `json:"core_password"`
}

// 注册处理函数
func Register(c *gin.Context) {
	const (
		EMPTY_USERNAME_OR_PASSWORD int = iota + 1
		USER_EXISTS
		FAILED_TO_CREATE_USER
	)

	errorMap := map[int]string{
		EMPTY_USERNAME_OR_PASSWORD: "Username or password cannot be empty",
		USER_EXISTS:                "Username already exists",
		FAILED_TO_CREATE_USER:      "Failed to create user"}

	var req RegisterRequestSerializer

	// 绑定JSON到User结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.CorePassword == "" {
		c.JSON(http.StatusOK, handler.NewErrorResponse(EMPTY_USERNAME_OR_PASSWORD, errorMap[EMPTY_USERNAME_OR_PASSWORD]))
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := dao.DB_INSTANCE.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(USER_EXISTS, errorMap[USER_EXISTS]))
		return
	}
	var user = model.User{Username: req.Username,
		Password:     req.Password,
		CorePassword: req.CorePassword}

	// 保存用户到数据库
	if err := dao.DB_INSTANCE.Create(&user).Error; err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(FAILED_TO_CREATE_USER, errorMap[FAILED_TO_CREATE_USER]))
		return
	}

	c.JSON(http.StatusOK, handler.NewSuccessResponse(nil))
}
