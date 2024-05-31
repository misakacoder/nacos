package config

import (
	"github.com/gin-gonic/gin"
	"nacos/router"
)

func RegisterV2(engine *gin.Engine) {
	cs := engine.Group(router.ApiV2 + "/cs")
	{
		cs.GET("/config/searchDetail", queryConfig)
	}
}
