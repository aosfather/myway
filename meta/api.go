package meta
//api的描述
type Api struct {
	Url string //api的url
	Desc string //描述
    Method []HttpMethod//允许的访问方法
	server *ServerCluster
	ServerUrl string //服务对应url
}
