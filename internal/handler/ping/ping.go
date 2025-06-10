package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	// c.JSON(http.StatusOK, handler.NewSuccessResponse(nil))
	c.String(http.StatusOK, "pong")
	return
}
