package controller

import (
	"duriand/internal/dao"
	"duriand/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 登录处理函数
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 绑定JSON输入
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var user model.User
	if err := dao.DB_INSTANCE.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "username": user.Username})
}
