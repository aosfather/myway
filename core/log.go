package core

import "time"

type Logger interface {
}

// http 访问方法
type HttpMethod byte

const (
	HM_GET  = 10
	HM_POST = 11
	HM_DEL  = 12
)

//http 状态
type HttpStatus int

const (
	HS_200 = 200
	HS_404 = 404
)

//访问日志
type AccessContent struct {
	Remote    string     //记录访问网站的客户端地址
	TimeStart time.Time  //记录访问时间
	Request   string     //用户的http请求起始行信息
	BodySent  int64      //服务器发送给客户端的响应body字节数
	Status    HttpStatus //http状态码，记录请求返回的状态码，例如：200、301、404等
	Referer   string     //记录此次请求是从哪个连接访问过来的，可以根据该参数进行防盗链设置。
	Url       string     //访问的URL地址
	Method    HttpMethod //访问的方法 GET POST DELETE
	Agent     string     //记录客户端访问信息，例如：浏览器、手机客户端等
	TimeEnd   time.Time  //记录访问结束时间
}

type ErrorContent struct {
	AccessContent
	ErrorCode    string //错误码
	ErrorMessage string //错误消息内容
	Descript     string //返回结果摘要
}
type AccessLogger interface {
	//服务访问
	ToAccess(content *AccessContent)
	//服务访问错误信息记录
	ToError(e *ErrorContent)

	WriteTextToAccess(text string)
}
