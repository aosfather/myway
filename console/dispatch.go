package console

import (
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

}