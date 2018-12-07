package main

import (
	"github.com/aosfather/myway/meta"
	"github.com/aosfather/myway/runtime"
)

func main() {
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
	server := meta.Server{}
	server.ID = 100
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
