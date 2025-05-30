package api

import (
	"duriand/internal/controller"
	"duriand/internal/dao"
	"duriand/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type InsertAccountRequestSerializer struct {
	Website  string `json:"website" binding:"required"`
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

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

	var req InsertAccountRequestSerializer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, controller.NewErrorResponse(INVALID_REQUEST, errorMap[INVALID_REQUEST]))
		return
	}

	// Get uid from JWT token
	uid := c.GetUint64("uid")
	// Get current timestamp
	now := time.Now().Unix()

	// Create account record
	account := model.Account{
		Uid:        uid,
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
		c.JSON(http.StatusOK, controller.NewErrorResponse(FAILED_TO_CREATE_ACCOUNT, errorMap[FAILED_TO_CREATE_ACCOUNT]))
		return
	}

	// Update timestamp using the user's uid
	timestamp := model.Timestamp{
		Uid:              uid,
		LatestUpdateTime: now,
	}

	if err := tx.Model(&model.Timestamp{}).Where("uid = ?", uid).
		Assign(map[string]interface{}{"latest_update_time": now}).
		FirstOrCreate(&timestamp).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, controller.NewErrorResponse(FAILED_TO_CREATE_ACCOUNT, errorMap[FAILED_TO_CREATE_ACCOUNT]))
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, controller.NewErrorResponse(FAILED_TO_CREATE_ACCOUNT, errorMap[FAILED_TO_CREATE_ACCOUNT]))
		return
	}

	c.JSON(http.StatusOK, controller.NewSuccessResponse("Account created successfully"))
}

type QueryAccountRequestSerializer struct {
	UpdateTime json.Number `json:"update_time" binding:"required"`
}

type PullMode int

const (
	PULL_ALL PullMode = iota + 1
	PULL_UPDATED
	PULL_NOTHING
)

type AccountRequestSerializer struct {
	Rid      uint64 `json:"rid"`
	Website  string `json:"website"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type QueryAccountResponseSerializer struct {
	PullMode   PullMode                   `json:"pull_mode"`
	UpdateTime int64                      `json:"update_time"`
	Accounts   []AccountRequestSerializer `json:"accounts"`
}

func QueryAccount(c *gin.Context) {
	const (
		NO_UPDATE_NEEDED int = iota + 1
		FAILED_TO_GET_TIMESTAMP
		FAILED_TO_GET_ACCOUNTS
	)

	errorMap := map[int]string{
		NO_UPDATE_NEEDED:        "No update needed",
		FAILED_TO_GET_TIMESTAMP: "Failed to get timestamp",
		FAILED_TO_GET_ACCOUNTS:  "Failed to get accounts",
	}

	// Get uid from JWT token
	uid := c.GetUint64("uid")

	// Parse request
	var req QueryAccountRequestSerializer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, controller.NewErrorResponse(FAILED_TO_GET_TIMESTAMP, errorMap[FAILED_TO_GET_TIMESTAMP]))
		return
	}

	// Get timestamp info
	var timestamp model.Timestamp
	if err := dao.DB_INSTANCE.Where("uid = ?", uid).First(&timestamp).Error; err != nil {
		c.JSON(http.StatusOK, controller.NewErrorResponse(FAILED_TO_GET_TIMESTAMP, errorMap[FAILED_TO_GET_TIMESTAMP]))
		return
	}

	response := QueryAccountResponseSerializer{
		UpdateTime: timestamp.LatestUpdateTime,
	}

	// Check update time
	var reqUpdateTime int64
	if ut, err := req.UpdateTime.Int64(); err != nil {
		c.JSON(http.StatusOK, controller.NewErrorResponse(FAILED_TO_GET_TIMESTAMP, errorMap[FAILED_TO_GET_TIMESTAMP]))
		return
	} else {
		reqUpdateTime = ut
	}
	fmt.Println(reqUpdateTime, "  ", timestamp.LatestUpdateTime)
	if reqUpdateTime >= timestamp.LatestUpdateTime {
		response.PullMode = PULL_NOTHING
		c.JSON(http.StatusOK, controller.NewSuccessResponse(response))
		return
	}

	var accounts []model.Account
	if reqUpdateTime < timestamp.LatestDeleteTime {
		response.PullMode = PULL_ALL
		// Return all accounts
		if err := dao.DB_INSTANCE.Where("uid = ?", uid).Find(&accounts).Error; err != nil {
			c.JSON(http.StatusOK, controller.NewErrorResponse(FAILED_TO_GET_ACCOUNTS, errorMap[FAILED_TO_GET_ACCOUNTS]))
			return
		}
	} else {
		response.PullMode = PULL_UPDATED
		// Return only updated accounts
		if err := dao.DB_INSTANCE.Where("uid = ? AND update_time > ?", uid, req.UpdateTime).Find(&accounts).Error; err != nil {
			c.JSON(http.StatusOK, controller.NewErrorResponse(FAILED_TO_GET_ACCOUNTS, errorMap[FAILED_TO_GET_ACCOUNTS]))
			return
		}
	}

	// Convert accounts to response format
	var accountResponses []AccountRequestSerializer
	for _, acc := range accounts {
		accountResponses = append(accountResponses, AccountRequestSerializer{
			Rid:      acc.Rid,
			Website:  acc.Website,
			Account:  acc.Account,
			Password: acc.Password,
		})
	}
	response.Accounts = accountResponses

	c.JSON(http.StatusOK, controller.NewSuccessResponse(response))
}
