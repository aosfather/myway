package meta

//常量定义
type HttpMethod byte //http 访问的方式
const (
	HTTP_GET    HttpMethod = 1
	HTTP_POST   HttpMethod = 2
	HTTP_PUT    HttpMethod = 3
	HTTP_DELETE HttpMethod = 4
	HTTP_HEAD   HttpMethod = 5
)

type HttpStatus byte //返回码
const (
	HS_400 HttpStatus = 4
	HS_500 HttpStatus = 5
	HS_200 HttpStatus = 2
)

//环境
type Env struct {
	ID          byte   //环境ID
	Name        string //环境名称
	Description string //环境描述
	Domain      string //环境对应的域名，可选
}
