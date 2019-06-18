package main

import (
	"github.com/aosfather/myway/console"
	"github.com/aosfather/myway/extends"
	"github.com/aosfather/myway/runtime"
)

//token生成的插件
const _PLUGIN_TOKEN = "tokens"

func main() {
	///*
	//   1、读取配置文件
	//   2、根据指定的目录，加载api的定义
	//*/
	//e := yamlConfig{}
	//e.Load("config.yaml")
	//
	////初始化日志组件
	//logfactory := logrusFactory{}
	//logfactory.Init(e.Config)
	//runtime.SetAccessLogger(logfactory.GetAccessLogger())
	//runtime.Log("test 001")
	//
	////启动控制端
	//admin := console.ConsoleDispatch{}
	//admin.Init(8980)
	//sh := console.SystemHandle{}
	//sh.Init(&admin)
	//go admin.Start()
	//
	////启动服务
	//dispatch := &runtime.DispatchManager{}
	//dispatch.Init()
	//handle := console.RegisterHandle{}
	//handle.Init(&admin, dispatch)
	//
	//proxy := runtime.HttpProxy{}
	//proxy.Init(dispatch)
	//
	////从配置中加载cluster的定义
	//LoadClusterFromFile(e.Config.ServerPath, dispatch)
	////从api的定义中加载代理的api
	//LoadAPIFromFile(e.Config.ApiPath, dispatch)
	//
	//
	//accss := runtime.AccessIntercepter{}
	//accss.Init("", 0, 7200, nil)
	//proxy.AddIntercepter(&accss)
	//proxy.Start()

	app := application{}
	app.Start()

}

type application struct {
	admin    console.ConsoleDispatch
	dispatch *runtime.DispatchManager
	proxy    *runtime.HttpProxy
}

func (this *application) Start() {
	e := yamlConfig{}
	e.Load("config.yaml")

	this.initLog(e.Config)
	//构建核心组件
	this.init()

	//加载管理端插件
	this.initAdminHandle()

	//加载api定义
	this.loadApisConfig(e.Config)

	//加载插件
	this.loadPlugins(e.Config, e.System)

	//启动控制端
	go this.admin.Start()
	//启动代理
	this.proxy.Start()

}

//构建核心组件
func (this *application) init() {
	//创建服务
	this.dispatch = &runtime.DispatchManager{}
	this.dispatch.Init()

	this.proxy = &runtime.HttpProxy{}
	this.proxy.Init(this.dispatch)

	//创建控制端
	this.admin = console.ConsoleDispatch{}
	this.admin.Init(8980)

}

func (this *application) loadApisConfig(config ApplicationConfigurate) {
	//从配置中加载cluster的定义
	LoadClusterFromFile(config.ServerPath, this.dispatch)
	//从api的定义中加载代理的api
	LoadAPIFromFile(config.ApiPath, this.dispatch)

}

func (this *application) initAdminHandle() {
	//增加系统接口
	sh := console.SystemHandle{}
	sh.Init(&this.admin)
	//增加注册接口
	handle := console.RegisterHandle{}
	handle.Init(&this.admin, this.dispatch)
}

func (this *application) initLog(config ApplicationConfigurate) {
	//初始化日志组件
	logfactory := logrusFactory{}
	logfactory.Init(config)
	runtime.SetAccessLogger(logfactory.GetAccessLogger())
	runtime.Log("test 001")
}

/**
  加载插件
*/
func (this *application) loadPlugins(config ApplicationConfigurate, system SystemConfigurate) {
	//加载权限配置
	rm := yamlAuthMetaManager{}
	rm.Load(config.RolePath)
	//加载用户配置
	um := userManager{}
	um.Load(config.UserPath)
	//加载access_token插件，token拦截及token生成的实现
	accss := extends.AccessTokenImp{}
	accss.Init(system.AuthRedis.Address, system.AuthRedis.DataBase, system.AuthRedis.Expire, &rm)
	accss.Add(&um)

	this.proxy.AddIntercepter(&accss)

	this.proxy.AddPlugin(_PLUGIN_TOKEN, accss.Call)
}
