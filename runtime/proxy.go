package runtime

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

//基本的代理
type HttpProxy struct {
	port   int
	server *fasthttp.Server
}

func (this *HttpProxy) Start() {
	this.server = &fasthttp.Server{Handler: this.ServeHTTP}
	addr := fmt.Sprintf("0.0.0.0:%d", 8080)
	this.server.ListenAndServe(addr)

}

func (this *HttpProxy) ServeHTTP(ctx *fasthttp.RequestCtx) {

}
