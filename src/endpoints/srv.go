package endpoints

import (
	"net"
	"time"
	"torch/src/structs"

	"github.com/gin-gonic/gin"
)

func srv(host string) (*structs.Srv, error) {
	_, addrs, err := net.LookupSRV("minecraft", "tcp", host)

	if err != nil {
		return nil, err
	}

	if len(addrs) < 1 {
		return nil, nil
	}

	return &structs.Srv{
		Target:     addrs[0].Target,
		Port:       addrs[0].Port,
		ObtainedAt: time.Now(),
		ExpiresAt:  time.Now().Add(time.Duration(srvCacheTime)),
	}, nil
}

func SrvHandler(c *gin.Context) {
	host := c.Param("host")

	data, err := srvCache.Value(host)
	if err == nil {
		c.JSON(200, data.Data().(structs.Srv))
		return
	}

	srv, err := srv(host)
	if err != nil || srv == nil {
		srv = &structs.Srv{
			Target:     host,
			Port:       25565,
			ObtainedAt: time.Now(),
			ExpiresAt:  time.Now().Add(time.Duration(srvCacheTime)),
		}
	}

	srvCache.Add(host, srvCacheTime, srv)
	c.JSON(200, srv)
}
