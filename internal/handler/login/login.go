package login

import (
	"duriand/internal/handler"
	login_service "duriand/internal/service/login"

	"net/http"

	"github.com/gin-gonic/gin"
)

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

	token, err := login_service.LoginService(req.Username, req.Password, req.CorePassword)
	if err != nil {
		c.JSON(http.StatusOK, handler.NewErrorResponse(INVALID_CREDENTIALS, errorMap[INVALID_CREDENTIALS]))
		return
	}

	c.JSON(http.StatusOK, handler.NewSuccessResponse(map[string]string{
		"token": token,
	}))
}
