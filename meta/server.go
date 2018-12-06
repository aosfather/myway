package meta

import "fmt"

//虚拟服务器
type ServerCluster struct {
	ID      string      //集群ID
	Name    string      //集群名称
	Servers []*Server   //服务
	Balance LoadBalance //负载策略
	Heath   HeathCheck  //健康检查
}

type Server struct {
	ID     int64  //服务器编号
	Ip     string //ipname
	Port   int    //端口
	MaxQPS int64  //支持的最大QPS
}

func (this *Server) Addr() string {
	return fmt.Sprintf("%s:%d", this.Ip, this.Port)
}

//健康监测
type HeathCheck struct {
	Path          string `json:"path"`
	Body          string `json:"body"`
	CheckInterval int64  `json:"checkInterval"`
	Timeout       int64  `json:"timeout"`
}

//重试策略
// RetryStrategy retry strategy
type RetryStrategy struct {
	Interval int32   `protobuf:"varint,1,opt,name=interval" json:"interval"`
	MaxTimes int32   `protobuf:"varint,2,opt,name=maxTimes" json:"maxTimes"`
	Codes    []int32 `protobuf:"varint,3,rep,name=codes" json:"codes,omitempty"`
}
