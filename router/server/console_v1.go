package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"nacos/configuration"
	"nacos/consts"
	"nacos/router"
	"nacos/router/auth"
	"nacos/router/cluster"
	"nacos/util"
	"net/http"
)

func RegisterV1(engine *gin.Engine) {
	server := engine.Group(router.ApiV1 + "/console/server")
	{
		server.GET("/state", serverState)
	}
	server.Use(auth.Auth)
	{
		server.GET("/guide", serverGuide)
		server.GET("/announcement", serverAnnouncement)
	}
}

func serverState(context *gin.Context) {
	nacos := configuration.Configuration.Nacos
	data := map[string]any{
		"defaultMaxSize":                "102400",
		"auth_system_type":              "nacos",
		"auth_enabled":                  fmt.Sprintf("%v", nacos.Auth.Enabled),
		"defaultMaxAggrSize":            "1024",
		"maxHealthCheckFailCount":       "12",
		"maxContent":                    "10485760",
		"console_ui_enabled":            "true",
		"defaultMaxAggrCount":           "10000",
		"defaultGroupQuota":             "200",
		"startup_mode":                  util.ConditionalExpression(len(cluster.Cluster.Clients) == 0, "standalone", "cluster"),
		"isHealthCheck":                 "true",
		"version":                       nacos.Version,
		"function_mode":                 nil,
		"isManageCapacity":              "true",
		"isCapacityLimitCheck":          "false",
		"datasource_platform":           "mysql",
		"notifyConnectTimeout":          "100",
		"server_port":                   "8848",
		"notifySocketTimeout":           "200",
		"defaultClusterQuota":           "100000",
		"login_page_enabled":            fmt.Sprintf("%v", nacos.Auth.Enabled),
		"plugin_datasource_log_enabled": "false",
	}
	context.JSON(http.StatusOK, data)
}

func serverGuide(context *gin.Context) {
	guide := `当前节点已关闭Nacos开源控制台使用，请修改application.properties中的nacos.console.ui.enabled参数为true打开开源控制台使用，详情查看<a href="https://nacos.io/zh-cn/docs/v2/guide/admin/console-guide.html">文档</a>中关于<code>关闭默认控制台部分</code>。`
	router.Success.With(guide).Ok(context)
}

func serverAnnouncement(context *gin.Context) {
	language := context.Query("language")
	router.Success.With(consts.GetI18n(language, consts.Announcement)).Ok(context)
}
