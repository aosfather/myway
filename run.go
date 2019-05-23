package main

import (
	"fmt"
	"github.com/aosfather/myway/console"
	"github.com/aosfather/myway/runtime"
)

func main() {
	e := yamlConfig{}
	e.Load("config.yaml")
	fmt.Println(e)

	//初始化日志组件
	logfactory := logrusFactory{}
	logfactory.Init(e.Config)
	runtime.SetAccessLogger(logfactory.GetAccessLogger())
	runtime.Log("test 001")

	//启动控制端
	admin := console.ConsoleDispatch{}
	admin.Init(8980)
	go admin.Start()

	//启动服务
	dispatch := &runtime.DispatchManager{}
	dispatch.Init()
	handle := console.RegisterHandle{}
	handle.Init(&admin, dispatch)

	////添加测试的api url，server
	//api := meta.Api{}
	//api.ServerUrl = "bb"
	//api.NameSpace = "m"
	//api.Url = "/my/a"
	//api.MaxQPS = 2
	//
	//cluster := meta.ServerCluster{}
	//cluster.ID = "test1"
	//cluster.Name = "test1"
	//cluster.Balance = 2
	//cluster.BalanceConfig = "test"
	//server := meta.Server{}
	//server.ID = 100
	//server.Tag.Init("test,dev")
	//server.Ip = "127.0.0.1"
	//server.Port = 8990
	//cluster.Servers = append(cluster.Servers, &server)
	//api.Cluster = &cluster
	//dispatch.AddCluster(&cluster)
	//dispatch.AddApi("", "", &api)

	//从配置中加载cluster的定义
	LoadClusterFromFile(e.Config.ServerPath, dispatch)
	//从api的定义中加载代理的api
	LoadAPIFromFile(e.Config.ApiPath, dispatch)

	proxy := runtime.HttpProxy{}
	proxy.Init(dispatch)
	proxy.Start()

}

/*
 TODO 系统加载次序
  1、读取配置文件
  2、根据指定的目录，加载api的定义
*/
