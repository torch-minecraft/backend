package endpoints

import (
	"fmt"
	"strconv"
	"strings"
	"torch/src/structs"

	"github.com/gin-gonic/gin"
)

func IconHandler(c *gin.Context) {
	ip := c.Param("ip")
	var port int

	if strings.Contains(ip, ":") {
		split := strings.Split(ip, ":")
		ip = split[0]
		p, err := strconv.Atoi(split[1])
		if err != nil {
			port = 25565
		}
		port = p
	} else {
		port = 25565
	}

	uintPort := uint16(port)

	cacheKey := fmt.Sprintf("%s:%d", ip, port)
	data, err := iconCache.Value(cacheKey)
	if err == nil {
		c.JSON(200, data.Data().(*structs.Icon))
		return
	}

	javaStatus, err := FetchJava(ip, uintPort)
	if err != nil {
		c.JSON(200, structs.Icon{
			Host: ip,
			Port: uintPort,
			Data: defaultIcon,
		})
		return
	}

	icon := structs.Icon{
		Host: ip,
		Port: uintPort,
		Data: javaStatus.Icon,
	}

	iconCache.Add(ip, iconCacheTime, icon)
	c.JSON(200, icon)
}
