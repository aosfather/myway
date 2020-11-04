package console

import (
	"encoding/json"
	"fmt"
	"github.com/aosfather/myway/runtime"
	"github.com/valyala/fasthttp"
)

/**
  控制台服务入口
*/
const (
	_DEFAULT_PORT = 9990
)

type ConsoleDispatch struct {
	port     int
	server   *fasthttp.Server
	dispatch *runtime.SimpleDispatch
}

func (this *ConsoleDispatch) Init(p int) {
	this.port = p
	this.dispatch = &runtime.SimpleDispatch{}
	this.dispatch.Init()

}

func (this *ConsoleDispatch) RegisterHandle(url string, h ConsoleHandler) {
	if url != "" && h != nil {
		this.dispatch.AddUrl(url, h)
	}

}

func (this *ConsoleDispatch) Start() {
	this.server = &fasthttp.Server{Handler: this.ServeHTTP}
	if this.port == 0 {
		this.port = _DEFAULT_PORT
	}
	addr := fmt.Sprintf("0.0.0.0:%d", this.port)
	this.server.ListenAndServe(addr)

}

func (this *ConsoleDispatch) ServeHTTP(ctx *fasthttp.RequestCtx) {
	//获取访问的url
	url := string(ctx.Request.URI().RequestURI())
	//通过dispatch，获取api的定义
	api := this.dispatch.GetUrl(url)

	if api == nil { //不存在的时候的处理
		ctx.Response.SetBodyString("the url not found!")
		return
	}

	//检查访问token

	//执行对应请求处理
	h := api.(ConsoleHandler)
	response := h(nil)

	ctx.WriteString(fmt.Sprintf("%d,%s,%s", response.Status, response.Code, response.Msg))

}

//控制请求
type ConsoleRequest struct {
	Token   string          `json:"access_token"`
	Version string          `json:"version"` //版本号
	Data    json.RawMessage `json:"data"`
}

//控制结果
type ConsoleResponse struct {
	Status byte        `json:"status"`
	Code   string      `json:"code"`
	Msg    string      `json:"message"`
	Data   interface{} `json:"data"`
}

//处理的hanler
type ConsoleHandler func(request *ConsoleRequest) ConsoleResponse
