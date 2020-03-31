package console

import (
	"encoding/json"
	"github.com/aosfather/myway/core"
	"os"
)

/**
系统控制指令
 1、shutdown 关闭服务器
 2、reload   重新加载 ,可选参数 api(定义) all（所有配置）
 3、restart 查新启动服务
*/
type reloadRequest struct {
	Tag string `json:"tag"`
}

type SystemHandle struct {
	App core.Application
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
	req := reloadRequest{}
	data, _ := request.Data.MarshalJSON()
	json.Unmarshal(data, &req)
	if this.App != nil {
		if req.Tag == "api" {
			this.App.ReloadApis()
			return ConsoleResponse{1, "001", "reload apis config success"}
		} else if req.Tag == "auth" {
			this.App.ReloadAuth()
			return ConsoleResponse{1, "001", "reload auth config success"}
		}
		//默认加载config
		this.App.ReloadConfig()
	}
	return ConsoleResponse{1, "001", "reload config success"}
}

func (this *SystemHandle) restart(request *ConsoleRequest) ConsoleResponse {
	return ConsoleResponse{}

}
