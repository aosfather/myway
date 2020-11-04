package meta

import (
	"fmt"
	"strings"
)

//应用映射
type ApplicationMapper struct {
	App     string         //应用名称
	Domain  string         //api对应的域名
	Host    string         //对应主机名称
	Cluster *ServerCluster //集群
	apis    []*ApiMapper
}

func (this *ApplicationMapper) AddMapper(api *ApiMapper) {
	if api != nil {
		api.application = this
		this.apis = append(this.apis, api)
	}
}

func (this *ApplicationMapper) GetMappers() []*ApiMapper {
	return this.apis
}

//API接口映射
type ApiMapper struct {
	application *ApplicationMapper
	Url         string //api的url
	TargetUrl   string //目标地址
	Label       string //接口标题
}

func (this *ApiMapper) GetCluster() *ServerCluster {
	if this.application != nil {
		return this.application.Cluster
	}
	return nil
}

func (this *ApiMapper) Key() string {
	return this.application.Domain + "/" + this.Url
}

func (this *ApiMapper) GetHost() string {
	return this.application.Host
}

//虚拟服务器
type ServerCluster struct {
	ID            string      //集群ID
	Name          string      //集群名称
	Servers       []*Server   //服务
	Balance       LoadBalance //负载策略
	BalanceConfig string      //配置
	Heath         HeathCheck  //健康检查
}

func (this *ServerCluster) AddServer(s *Server) {
	if s != nil {
		this.Servers = append(this.Servers, s)
	}
}

type Server struct {
	ID     int64  //服务器编号
	Ip     string //ipname
	Port   int    //端口
	MaxQPS int64  //支持的最大QPS
	Tag    Tags
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

type Tags struct {
	tags []string
}

func (this *Tags) Init(str string) {
	this.tags = strings.Split(str, ",")
}

func (this *Tags) Has(tag string) bool {
	for _, v := range this.tags {
		if v == tag {
			return true
		}
	}

	return false
}
