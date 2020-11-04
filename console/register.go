package console

import (
	"encoding/json"
	//"github.com/aosfather/myway/meta"
	"github.com/aosfather/myway/runtime"
	"sync"
)

//注册接口
type RegisterHandle struct {
	dm    *runtime.DispatchManager
	mutex sync.RWMutex
}

type clusterInfo struct {
	ID            string //集群ID
	Name          string //集群名称
	Balance       int32  //负载策略
	BalanceConfig string //配置
}

func (this *RegisterHandle) Init(c *ConsoleDispatch, r *runtime.DispatchManager) {
	this.dm = r
	c.RegisterHandle("/meta/addcluster", this.addCluster)
	c.RegisterHandle("/meta/addserver", this.addServer)
	c.RegisterHandle("/meta/addapi", this.addApi)

}

//新增集群
func (this *RegisterHandle) addCluster(request *ConsoleRequest) ConsoleResponse {
	info := clusterInfo{}
	json.Unmarshal(request.Data, &info)
	//cluster := this.dm.GetCluster(info.Name)
	//if cluster != nil {
	//	return ConsoleResponse{}
	//}
	//this.mutex.Lock()
	//c := &meta.ServerCluster{}
	//c.ID = info.ID
	//c.Name = info.Name
	//c.Balance = meta.LoadBalance(info.Balance)
	//c.BalanceConfig = info.BalanceConfig
	//this.dm.AddCluster(c)
	this.mutex.Unlock()

	return ConsoleResponse{}
}

type serverInfo struct {
	ID      int64  `json:"id"`
	Tag     string `json:"tag"` //tag标签 a,b格式
	Ip      string `json:"ip"`
	Port    int    `json:"port"`
	Cluster string `json:"cluster"` //集群名称
}

//新增服务器
func (this *RegisterHandle) addServer(request *ConsoleRequest) ConsoleResponse {
	info := serverInfo{}
	json.Unmarshal(request.Data, &info)
	this.mutex.Lock()
	//cluster := this.dm.GetCluster(info.Cluster)
	//if cluster != nil {
	//	s := &meta.Server{}
	//	s.ID = info.ID
	//	s.Tag.Init(info.Tag)
	//	s.Port = info.Port
	//	s.Ip = info.Ip
	//
	//	cluster.Servers = append(cluster.Servers, s)
	//} else {
	//	return ConsoleResponse{100, "404", "cluster not found!"}
	//}
	defer this.mutex.Unlock()

	return ConsoleResponse{200, "200", "success!", nil}
}

type apiInfo struct {
	Cluster   string `json:"cluster"` //集群
	Url       string `json:"api"`
	NameSpace string `json:"ns"`
	ServerUrl string `json:"url"`
	MaxQPS    int64  `json:"qps"`
}

//批量注册api
func (this *RegisterHandle) addApi(request *ConsoleRequest) ConsoleResponse {
	if request == nil {
		return ConsoleResponse{100, "500", "the request is nil", nil}
	}
	info := &apiInfo{}
	json.Unmarshal(request.Data, info)
	this.mutex.Lock()
	//c := this.dm.GetCluster(info.Cluster)
	//if c != nil {
	//	//api := meta.Api{}
	//	//api.Url = info.Url
	//	//api.NameSpace = info.NameSpace
	//	//api.ServerUrl = info.ServerUrl
	//	//api.MaxQPS = info.MaxQPS
	//	//api.Cluster = c
	//	//this.dm.AddApi("", "", &api)
	//
	//} else {
	//	return ConsoleResponse{100, "404", "cluster not found!"}
	//}
	return ConsoleResponse{200, "200", "success!", nil}
}
