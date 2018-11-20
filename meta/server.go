package meta
//虚拟服务器
type ServerCluster struct {
    Name string
    servers []*Server
}

type Server struct {
	Ip string  //ipname
	Port int   //端口
}
