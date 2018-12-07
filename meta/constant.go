package meta

//常量定义
type HttpMethod byte //http 访问的方式
const (
	HTTP_GET  HttpMethod = 1
	HTTP_POST HttpMethod = 2
)

//环境
type Env struct {
	ID          byte   //环境ID
	Name        string //环境名称
	Description string //环境描述
	Domain      string //环境对应的域名，可选
}
