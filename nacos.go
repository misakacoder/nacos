package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/misakacoder/logger"
	"nacos/configuration"
	"nacos/consts"
	"nacos/middleware"
	"nacos/model"
	"nacos/router"
	"nacos/router/auth"
	"nacos/router/cluster"
	"nacos/router/config"
	"nacos/router/namespace"
	"nacos/router/server"
	"nacos/util"
	"strings"
	"time"
)

//go:embed ui
var ui embed.FS
var startTime time.Time

var (
	version   = "dev"
	buildTime = "unknown"
)

func init() {
	configuration.Configuration.Nacos.Version = version
	startTime = time.Now()
	gin.SetMode(gin.ReleaseMode)
	initLogger()
}

func main() {
	engine := gin.New()
	initTable()
	initEngine(engine)
	startup(engine)
}

func initLogger() {
	log := configuration.Configuration.Log
	logger.SetLogger(logger.NewSimpleLogger(log.Filename))
	level, _ := logger.Parse(log.Level)
	logger.SetLevel(level)
}

func initTable() {
	model.FirstOrCreate()
}

func initEngine(engine *gin.Engine) {
	engine.Use(middleware.FS("", middleware.EmbedFile(ui)))
	engine.Use(middleware.NetWork)
	engine.Use(middleware.CSRF)
	engine.Use(middleware.Recovery)
	engine.NoRoute(router.NotFound)
	config.RegisterV1(engine)
	config.RegisterV2(engine)
	server.RegisterV1(engine)
	namespace.RegisterV1(engine)
	auth.RegisterV1(engine)
	cluster.RegisterV1(engine)
}

func startup(engine *gin.Engine) {
	logger.Info("The Nacos version is %s and the build time is %s", version, buildTime)
	conf := configuration.Configuration.Server
	bind := conf.Bind
	port := conf.Port
	banner := strings.Builder{}
	startUpTime := time.Since(startTime)
	banner.WriteString(fmt.Sprintf("Started Nacos in %.2f seconds...", startUpTime.Seconds()))
	addresses := util.ConditionalExpression(bind == consts.AnyAddress, util.GetLocalAddr(), []string{bind})
	for _, address := range addresses {
		banner.WriteString(fmt.Sprintf("\n - Listen on: http://%s:%d", address, port))
	}
	logger.Info(banner.String())
	err := engine.Run(fmt.Sprintf("%s:%d", bind, port))
	if err != nil {
		logger.Error(err.Error())
	}
}
