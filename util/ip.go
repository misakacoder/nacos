package util

import (
	"github.com/gin-gonic/gin"
	"github.com/misakacoder/logger"
	"nacos/consts"
	"net"
)

func GetLocalAddr() (addr []string) {
	addr = append(addr, consts.Localhost)
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, value := range interfaces {
		if (value.Flags & net.FlagUp) != 0 {
			addresses, _ := value.Addrs()
			for _, address := range addresses {
				if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						addr = append(addr, ipNet.IP.String())
					}
				}
			}
		}
	}
	return
}

func GetClientIP(context *gin.Context) string {
	ip := context.ClientIP()
	return ConditionalExpression(ip == "::1", consts.Localhost, ip)
}
