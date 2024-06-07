package cluster

import (
	"fmt"
	"nacos/configuration"
	"nacos/consts"
	"nacos/util"
	"net/rpc"
	"time"
)

var CLUSTER = NewCluster()

type Args struct {
	Token string
	Data  any
}

type Node struct {
	Address     string
	State       string
	RefreshTime time.Time
	client      *rpc.Client
}

type Cluster struct {
	Master    *Node
	token     string
	rpcServer *rpc.Server
	Slaves    []*Node
}

func NewCluster() *Cluster {
	server := configuration.Configuration.Server
	ip := util.ConditionalExpression(server.Bind == consts.AnyAddress, consts.Localhost, server.Bind)
	conf := configuration.Configuration.Nacos.Cluster
	addresses := conf.List
	cluster := &Cluster{
		Master:    &Node{Address: fmt.Sprintf("%s:%d", ip, server.Port), State: "UP", RefreshTime: time.Now()},
		token:     conf.Token,
		rpcServer: rpc.NewServer(),
	}
	cluster.rpcServer.Register(cluster)
	cluster.rpcServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	if len(addresses) > 0 {
		for _, address := range addresses {
			cluster.Slaves = append(cluster.Slaves, &Node{Address: address, State: "DOWN", RefreshTime: time.Now(), client: nil})
		}
		go func(cluster *Cluster) {
			args := Args{Token: conf.Token}
			for {
				for _, element := range cluster.Slaves {
					element.RefreshTime = time.Now()
					if element.client == nil {
						addr := element.Address
						client, err := rpc.DialHTTP("tcp", addr)
						if err == nil {
							element.client = client
						} else {
							continue
						}
					}
					ok := false
					err := element.client.Call("Cluster.Heartbeat", args, &ok)
					if err != nil {
						element.State = "DOWN"
						element.client = nil
					} else if ok {
						element.State = "UP"
					} else {
						element.State = "SUSPICIOUS"
					}
				}
				time.Sleep(5 * time.Second)
			}
		}(cluster)
	}
	return cluster
}

func (cluster *Cluster) Heartbeat(args *Args, ok *bool) error {
	if args.Token == cluster.token {
		*ok = true
	}
	return nil
}
