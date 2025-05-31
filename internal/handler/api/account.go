package api

import (
	"duriand/internal/handler"
	service_api "duriand/internal/service/api"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InsertAccountHandler(c *gin.Context) {
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

	type InsertAccountRequestSerializer struct {
		Website  string `json:"website" binding:"required"`
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req InsertAccountRequestSerializer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(INVALID_REQUEST, errorMap[INVALID_REQUEST]))
		return
	}

	// Get uid from JWT token
	uid := c.GetUint64("uid")
	if err := service_api.InsertAccountService(uid, req.Website, req.Account, req.Password); err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(FAILED_TO_CREATE_ACCOUNT, errorMap[FAILED_TO_CREATE_ACCOUNT]+": "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, handler.NewSuccessResponse(nil))
}

func QueryAccountHandler(c *gin.Context) {
	type QueryAccountRequestSerializer struct {
		UpdateTime json.Number `json:"update_time" binding:"required"`
	}

	type QueryAccountResponseSerializer struct {
		PullMode   service_api.PullMode         `json:"pull_mode"`
		UpdateTime int64                        `json:"update_time"`
		Accounts   []service_api.QueriedAccount `json:"accounts"`
	}

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

	// Parse request
	var req QueryAccountRequestSerializer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(FAILED_TO_GET_TIMESTAMP, errorMap[FAILED_TO_GET_TIMESTAMP]))
		return
	}
	// Get uid from JWT token
	uid := c.GetUint64("uid")
	updateTime, err := req.UpdateTime.Int64()
	if err != nil || updateTime < 0 {
		c.JSON(http.StatusOK, handler.NewErrorResponse(FAILED_TO_GET_TIMESTAMP, errorMap[FAILED_TO_GET_TIMESTAMP]))
		return
	}
	// call servie api
	pullMode, latestUpdateTime, queriedAccounts, err := service_api.QueryAccountService(uid, updateTime)
	if err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(FAILED_TO_GET_ACCOUNTS, errorMap[FAILED_TO_GET_ACCOUNTS]+": "+err.Error()))
		return
	}
	response := QueryAccountResponseSerializer{
		PullMode:   pullMode,
		UpdateTime: latestUpdateTime,
		Accounts:   queriedAccounts,
	}

	c.JSON(http.StatusOK, handler.NewSuccessResponse(response))
}

func UpdateAccountHandler(c *gin.Context) {
	const (
		INVALID_REQUEST int = iota + 1
		FAILED_TO_UPDATE_ACCOUNT
	)

	errorMap := map[int]string{
		INVALID_REQUEST:          "Invalid request data",
		FAILED_TO_UPDATE_ACCOUNT: "Failed to update account record",
	}

	type UpdateAccountRequestSerializer struct {
		Rid      json.Number `json:"rid" binding:"required"`
		Website  string      `json:"website"`
		Account  string      `json:"account"`
		Password string      `json:"password"`
	}

	var req UpdateAccountRequestSerializer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(INVALID_REQUEST, errorMap[INVALID_REQUEST]))
		return
	}

	uid := c.GetUint64("uid")
	var rid uint64
	if ridi64, err := req.Rid.Int64(); err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(INVALID_REQUEST, errorMap[INVALID_REQUEST]))
		return
	} else {
		rid = uint64(ridi64)
	}
	// call service api
	if err := service_api.UpdateAccountService(uid, rid, req.Website, req.Account, req.Password); err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(FAILED_TO_UPDATE_ACCOUNT, errorMap[FAILED_TO_UPDATE_ACCOUNT]+": "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, handler.NewSuccessResponse(nil))
}

func DeleteAccountHandler(c *gin.Context) {
	const (
		INVALID_REQUEST int = iota + 1
		FAILED_TO_DELETE_ACCOUNT
	)

	errorMap := map[int]string{
		INVALID_REQUEST:          "Invalid request data",
		FAILED_TO_DELETE_ACCOUNT: "Failed to delete account record",
	}

	type DeleteAccountRequestSerializer struct {
		Rid json.Number `json:"rid" binding:"required"`
	}

	var req DeleteAccountRequestSerializer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(INVALID_REQUEST, errorMap[INVALID_REQUEST]))
		return
	}

	uid := c.GetUint64("uid")
	var rid uint64
	if ridi64, err := req.Rid.Int64(); err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(INVALID_REQUEST, errorMap[INVALID_REQUEST]))
		return
	} else {
		rid = uint64(ridi64)
	}
	if err := service_api.DeleteAccountService(uid, rid); err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(FAILED_TO_DELETE_ACCOUNT, errorMap[FAILED_TO_DELETE_ACCOUNT]+": "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, handler.NewSuccessResponse(nil))
}
