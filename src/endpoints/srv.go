package endpoints

import (
	"net"

	"github.com/gin-gonic/gin"
)

func LookupSrv(host string) (*net.SRV, error) {
	_, addrs, err := net.LookupSRV("minecraft", "tcp", host)

	if err != nil {
		return nil, err
	}

	if len(addrs) < 1 {
		return nil, nil
	}

	return addrs[0], nil
}

func SrvHandler(c *gin.Context) {
	host := c.Query("host")

	data, err := srvCache.Value(host)
	if err == nil {
		c.JSON(200, data.Data().(*net.SRV))
		return
	}

	srv, err := LookupSrv(host)
	if err != nil || srv == nil {
		srv = &net.SRV{
			Target: host,
			Port:   25565,
		}
	}

	srvCache.Add(host, srvCacheTime, srv)
	c.JSON(200, srv)
}
