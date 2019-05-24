package main

import (
	"github.com/aosfather/myway/console"
	"github.com/aosfather/myway/runtime"
)

func main() {
	/*
	   1、读取配置文件
	   2、根据指定的目录，加载api的定义
	*/
	e := yamlConfig{}
	e.Load("config.yaml")

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

	//从配置中加载cluster的定义
	LoadClusterFromFile(e.Config.ServerPath, dispatch)
	//从api的定义中加载代理的api
	LoadAPIFromFile(e.Config.ApiPath, dispatch)

	proxy := runtime.HttpProxy{}
	proxy.Init(dispatch)
	proxy.Start()

}
