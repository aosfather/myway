package main

import (
	"github.com/aosfather/myway/console"
	"github.com/aosfather/myway/extends"
	"github.com/aosfather/myway/runtime"
)

//token生成的插件
const _PLUGIN_TOKEN = "tokens"

func main() {
	app := application{}
	app.Start()

}

//应用容器，会实现容器接口，用于做admin 接口的控制
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

//加载api的定义
func (this *application) loadApisConfig(config ApplicationConfigurate) {
	//从配置中加载cluster的定义
	LoadClusterFromFile(config.ServerPath, this.dispatch)
	//从api的定义中加载代理的api
	LoadAPIFromFile(config.ApiPath, this.dispatch)

}

//初始化管理控制台
func (this *application) initAdminHandle() {
	//增加系统接口
	sh := console.SystemHandle{}
	sh.Init(&this.admin)
	//增加注册接口
	handle := console.RegisterHandle{}
	handle.Init(&this.admin, this.dispatch)
}

//初始化日志
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

	this.proxy.AddIntercepter(&accss)               //token拦截
	this.proxy.AddPlugin(_PLUGIN_TOKEN, accss.Call) //token生成插件服务
}
