package console

import "os"

/**
系统控制指令
 1、shutdown 关闭服务器
 2、reload   重新加载 ,可选参数 api(定义) all（所有配置）
 3、restart 查新启动服务
*/

//应用容器
type ApplicationContent interface {
	ReloadConfig()   //重新加载配置
	ReloadAuth()     //重新加载权限
	Restart()        //重启
	ShutdownNotify() //关闭通知

}
type SystemHandle struct {
	App ApplicationContent
}

func (this *SystemHandle) Init(c *ConsoleDispatch) {
	c.RegisterHandle("/shutdown", this.shutdown)
	c.RegisterHandle("/reload", this.reload)
	c.RegisterHandle("/restart", this.restart)

}

func (this *SystemHandle) shutdown(request *ConsoleRequest) ConsoleResponse {
	if this.App != nil {
		this.App.ShutdownNotify()
	}
	defer os.Exit(0)
	return ConsoleResponse{1, "001", "exit by you!"}
}

func (this *SystemHandle) reload(request *ConsoleRequest) ConsoleResponse {
	if this.App != nil {
		this.App.ReloadConfig()
	}
	return ConsoleResponse{1, "001", "reload config success"}
}

func (this *SystemHandle) restart(request *ConsoleRequest) ConsoleResponse {
	return ConsoleResponse{}

}
