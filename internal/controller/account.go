package controller

import (
	"duriand/internal/dao"
	"duriand/internal/model"
	"duriand/internal/serializer"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func InsertAccount(c *gin.Context) {
	const (
		INVALID_REQUEST int = iota + 1
		FAILED_TO_CREATE_ACCOUNT
		FAILED_TO_GET_USER
	)

	errorMap := map[int]string{
		INVALID_REQUEST:          "Invalid request data",
		FAILED_TO_CREATE_ACCOUNT: "Failed to create account record",
		FAILED_TO_GET_USER:       "Failed to get user information",
	}

	var req serializer.InsertAccountRequestSerializer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, serializer.NewErrorResponse(INVALID_REQUEST, errorMap[INVALID_REQUEST]))
		return
	}

	// Get username from JWT token
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusOK, serializer.NewErrorResponse(FAILED_TO_GET_USER, errorMap[FAILED_TO_GET_USER]))
		return
	}

	// Get user info from database
	var user model.User
	if err := dao.DB_INSTANCE.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, serializer.NewErrorResponse(FAILED_TO_GET_USER, errorMap[FAILED_TO_GET_USER]))
		return
	}

	// Get current timestamp
	now := time.Now().Unix()

	// Create account record
	account := model.Account{
		Uid:        user.Uid,
		Website:    req.Website,
		Account:    req.Account,
		Password:   req.Password,
		UpdateTime: now,
	}

	// Start transaction
	tx := dao.DB_INSTANCE.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Insert account record
	if err := tx.Create(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, serializer.NewErrorResponse(FAILED_TO_CREATE_ACCOUNT, errorMap[FAILED_TO_CREATE_ACCOUNT]))
		return
	}

	// Update timestamp using the user's uid
	timestamp := model.Timestamp{
		Uid:              user.Uid,
		LatestUpdateTime: now,
	}

	if err := tx.Model(&model.Timestamp{}).Where("uid = ?", user.Uid).
		Assign(map[string]interface{}{"latest_update_time": now}).
		FirstOrCreate(&timestamp).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, serializer.NewErrorResponse(FAILED_TO_CREATE_ACCOUNT, errorMap[FAILED_TO_CREATE_ACCOUNT]))
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, serializer.NewErrorResponse(FAILED_TO_CREATE_ACCOUNT, errorMap[FAILED_TO_CREATE_ACCOUNT]))
		return
	}

	c.JSON(http.StatusOK, serializer.NewSuccessResponse("Account created successfully"))
}
