package config

import (
	"github.com/gin-gonic/gin"
	"nacos/router"
	"nacos/router/auth"
	"net/http"
)

func RegisterV1(engine *gin.Engine) {
	cs := engine.Group(router.ApiV1+"/cs", auth.Auth)
	{
		cs.GET("/configs", queryConfig)
		cs.GET("/searchDetail", queryConfig)
		cs.POST("/configs", saveConfig)
		cs.DELETE("/configs", delConfig)
	}
	{
		cs.GET("/history/configs", listHistoryConfig)
		cs.GET("/history", queryHistoryConfig)
	}
	{
		cs.GET("/configs/listener", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{"collectStatus": 200, "lisentersGroupkeyStatus": map[string]string{}})
		})
		cs.POST("/configs/listener", listenConfig)
	}
}
