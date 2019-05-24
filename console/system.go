package console

import "os"

/**
系统控制指令
 1、shutdown 关闭服务器
 2、reload   重新加载 ,可选参数 api(定义) all（所有配置）
 3、restart 查新启动服务
*/

type SystemHandle struct {
}

func (this *SystemHandle) Init(c *ConsoleDispatch) {
	c.RegisterHandle("/shutdown", this.shutdown)
	c.RegisterHandle("/reload", this.reload)
	c.RegisterHandle("/restart", this.restart)

}

func (this *SystemHandle) shutdown(request *ConsoleRequest) ConsoleResponse {
	defer os.Exit(0)
	return ConsoleResponse{1, "001", "exit by you!"}
}

func (this *SystemHandle) reload(request *ConsoleRequest) ConsoleResponse {
	return ConsoleResponse{}

}

func (this *SystemHandle) restart(request *ConsoleRequest) ConsoleResponse {
	return ConsoleResponse{}

}
