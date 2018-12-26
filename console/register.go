package console

import (
	"github.com/aosfather/myway/runtime"
	"sync"
)

//注册接口
type RegisterHandle struct {
	dm    *runtime.DispatchManager
	mutex sync.RWMutex
}

//新增集群
func (this *RegisterHandle) addCluster(request *ConsoleRequest) ConsoleResponse {
	cluster := this.dm.GetCluster("")
	if cluster != nil {
		return ConsoleResponse{}
	}
	this.mutex.Lock()
	this.dm.AddCluster(nil)
	this.mutex.Unlock()

	return ConsoleResponse{}
}

//新增服务器
func (this *RegisterHandle) addServer(request *ConsoleRequest) ConsoleResponse {
	this.mutex.Lock()

	defer this.mutex.Unlock()

	return ConsoleResponse{}
}

//批量注册api
func (this *RegisterHandle) addApi(request *ConsoleRequest) ConsoleResponse {

	return ConsoleResponse{}
}
