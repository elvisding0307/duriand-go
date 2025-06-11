package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyHandler(c *gin.Context) {
	// c.JSON(http.StatusOK, handler.NewSuccessResponse(nil))
	c.String(http.StatusOK, "")
	return
}
