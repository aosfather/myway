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

//常用的请求头
const (
	UserAgent               = "User-Agent"
	ContentType             = "Content-Type"
	ContentLength           = "Content-Length"
	Authorization           = "Authorization"
	ContentEncoding         = "Content-Encoding"
	Accept                  = "Accept"
	AcceptEncoding          = "Accept-Encoding"
	StrictTransportSecurity = "Strict-Transport-Security"
	CacheControl            = "Cache-Control"
	Pragma                  = "Pragma"
	Expires                 = "Expires"
	Connection              = "Connection"
	XRealIP                 = "X-Real-IP"
	XForwardFor             = "X-Forwarded-For"
)

//常用的content type
const (
	ApplicationJSON = "application/json"
	ApplicationXML  = "application/xml"
)
