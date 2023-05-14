package endpoints

import "github.com/gin-gonic/gin"

func PingHandler(c *gin.Context) {
	c.JSON(200, "pong")
}
