package register

import (
	"duriand/internal/handler"
	register_service "duriand/internal/service/register"
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

	if err := register_service.RegisterService(req.Username, req.Password, req.CorePassword); err != nil {
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
