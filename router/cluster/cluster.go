package cluster

import (
	"github.com/gin-gonic/gin"
	"nacos/cluster"
	"nacos/model"
	"nacos/router"
	"nacos/router/auth"
	"nacos/util"
	"net/http"
	"net/rpc"
	"strings"
	"time"
)

func RegisterV1(engine *gin.Engine) {
	engine.GET(rpc.DefaultDebugPath, func(context *gin.Context) {
		http.DefaultServeMux.ServeHTTP(context.Writer, context.Request)
	})
	engine.Handle(http.MethodConnect, rpc.DefaultRPCPath, func(context *gin.Context) {
		http.DefaultServeMux.ServeHTTP(context.Writer, context.Request)
	})

	cluster := engine.Group(router.ApiV1+"/core/cluster", auth.Auth)
	{
		cluster.GET("/nodes", nodes)
		cluster.POST("/server/leave", leave)
	}
}

func nodes(context *gin.Context) {
	keyword := context.Query("keyword")
	clients := []*cluster.Node{cluster.CLUSTER.Master}
	clients = append(clients, cluster.CLUSTER.Slaves...)
	node := make([]model.Node, 0)
	for _, client := range clients {
		if strings.HasPrefix(client.Address, keyword) {
			node = append(node, getNode(client))
		}
	}
	router.OK.With(node).Ok(context)
}

func leave(context *gin.Context) {
	router.AccessDenied.Error(context)
}

func getNode(client *cluster.Node) model.Node {
	ip, port := parseAddress(client.Address)
	node := model.Node{
		IP:      ip,
		Port:    port,
		State:   client.State,
		Address: client.Address,
		Metadata: map[string]any{
			"lastRefreshTime": client.RefreshTime.Format(time.DateTime),
		},
	}
	return node
}

func parseAddress(address string) (string, int) {
	addresses := strings.Split(address, ":")
	return addresses[0], util.Atoi[int](addresses[1])
}
