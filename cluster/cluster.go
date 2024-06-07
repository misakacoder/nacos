package cluster

import (
	"encoding/gob"
	"fmt"
	"github.com/misakacoder/logger"
	"nacos/configuration"
	"nacos/consts"
	"nacos/listener"
	"nacos/model"
	"nacos/util"
	"net/rpc"
	"time"
)

var (
	CLUSTER    = NewCluster()
	up         = "UP"
	down       = "DOWN"
	suspicious = "SUSPICIOUS"
	result     = &struct{}{}
)

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
	gob.Register(model.ConfigKey{})
	server := configuration.Configuration.Server
	ip := util.ConditionalExpression(server.Bind == consts.AnyAddress, consts.Localhost, server.Bind)
	conf := configuration.Configuration.Nacos.Cluster
	addresses := conf.List
	cluster := &Cluster{
		Master:    &Node{Address: fmt.Sprintf("%s:%d", ip, server.Port), State: up, RefreshTime: time.Now()},
		token:     conf.Token,
		rpcServer: rpc.NewServer(),
	}
	cluster.rpcServer.Register(cluster)
	cluster.rpcServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	if len(addresses) > 0 {
		for _, address := range addresses {
			cluster.Slaves = append(cluster.Slaves, &Node{Address: address, State: down, RefreshTime: time.Now(), client: nil})
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
						element.State = down
						element.client = nil
					} else if ok {
						element.State = up
					} else {
						element.State = suspicious
					}
				}
				time.Sleep(5 * time.Second)
			}
		}(cluster)
	}
	return cluster
}

func (cluster *Cluster) Heartbeat(args *Args, ok *bool) error {
	return cluster.auth(args, func() error {
		*ok = true
		return nil
	})
}

func (cluster *Cluster) NotifyConfigListener(args *Args, result *struct{}) error {
	return cluster.auth(args, func() error {
		key := args.Data.(model.ConfigKey)
		listener.ConfigListenerManager.Notify(key)
		return nil
	})
}

func (cluster *Cluster) NotifySlaveConfigListener(configKey model.ConfigKey) {
	for _, slave := range cluster.Slaves {
		if slave.State == up {
			slave.client.Call("Cluster.NotifyConfigListener", cluster.buildArgs(configKey), result)
		}
	}
}

func (cluster *Cluster) buildArgs(data any) *Args {
	return &Args{
		Token: cluster.token,
		Data:  data,
	}
}

func (cluster *Cluster) auth(args *Args, fn func() error) error {
	if args.Token == cluster.token {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("%v", util.GetStackTrace(err))
			}
		}()
		return fn()
	}
	return nil
}
