package main

import (
	"github.com/aosfather/myway/console"
	"github.com/aosfather/myway/meta"
	"github.com/aosfather/myway/runtime"
)

func main() {

	//启动控制端
	admin := console.ConsoleDispatch{}
	admin.Init(8980)
	go admin.Start()

	//启动服务
	dispatch := &runtime.DispatchManager{}
	dispatch.Init()

	//添加测试的api url，server
	api := meta.Api{}
	api.ServerUrl = "meta/query"
	api.NameSpace = "m"
	api.Url = "/a"
	api.MaxQPS = 2

	cluster := meta.ServerCluster{}
	cluster.ID = "test"
	cluster.Name = "测试集群"
	cluster.Balance = 2
	cluster.BalanceConfig = "test"
	server := meta.Server{}
	server.ID = 100
	server.Tag.Init("test,dev")
	server.Ip = "127.0.0.1"
	server.Port = 8990
	cluster.Servers = append(cluster.Servers, &server)
	api.Cluster = &cluster
	dispatch.AddCluster(&cluster)
	dispatch.AddApi("", "", &api)

	proxy := runtime.HttpProxy{}
	proxy.Init(dispatch)
	proxy.Start()

}
