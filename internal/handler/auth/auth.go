package auth

import (
	"duriand/internal/handler"
	service_auth "duriand/internal/service/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	type RegisterRequestSerializer struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		CorePassword string `json:"core_password"`
	}

	const (
		EMPTY_USERNAME_OR_PASSWORD int = iota + 1
		USER_EXISTS
		FAILED_TO_CREATE_USER
	)

	errorMap := map[int]string{
		EMPTY_USERNAME_OR_PASSWORD: "Username or password cannot be empty",
		USER_EXISTS:                "Username already exists",
		FAILED_TO_CREATE_USER:      "Failed to create user",
	}

	var req RegisterRequestSerializer

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.CorePassword == "" {
		c.JSON(http.StatusOK, handler.NewErrorResponse(EMPTY_USERNAME_OR_PASSWORD, errorMap[EMPTY_USERNAME_OR_PASSWORD]))
		return
	}

	if err := service_auth.RegisterService(req.Username, req.Password, req.CorePassword); err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusOK, handler.NewErrorResponse(USER_EXISTS, errorMap[USER_EXISTS]))
			return
		} else {
			c.JSON(http.StatusOK, handler.NewErrorResponse(FAILED_TO_CREATE_USER, errorMap[FAILED_TO_CREATE_USER]))
			return

		}
	}

	c.JSON(http.StatusOK, handler.NewSuccessResponse(nil))
}

func LoginHandler(c *gin.Context) {
	type LoginRequestSerializer struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		CorePassword string `json:"core_password"`
	}

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

	token, err := service_auth.LoginService(req.Username, req.Password, req.CorePassword)
	if err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(INVALID_CREDENTIALS, errorMap[INVALID_CREDENTIALS]))
		return
	}

	c.JSON(http.StatusOK, handler.NewSuccessResponse(map[string]string{
		"token": token,
	}))
}
