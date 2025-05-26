package controller

import (
	"duriand/internal/dao"
	"duriand/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 注册处理函数
func Register(c *gin.Context) {
	var requestBody struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		CorePassword string `json:"core_password"`
	}

	// 绑定JSON到User结构体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := dao.DB_INSTANCE.Where("username = ?", requestBody.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}
	var user = model.User{Username: requestBody.Username,
		Password:     requestBody.Password,
		CorePassword: requestBody.CorePassword}

	// 保存用户到数据库
	if err := dao.DB_INSTANCE.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
