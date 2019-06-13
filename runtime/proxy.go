package runtime

import (
	"fmt"
	"github.com/aosfather/myway/meta"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"time"
)

const Cluster_Static = "-" //静态资源
//基本的代理
type HttpProxy struct {
	port         int
	server       *fasthttp.Server
	dispatch     *DispatchManager //分发管理
	intercepters []Intercepter    //拦截器处理
	client       *FastHTTPClient
}

func (this *HttpProxy) Init(dispatch *DispatchManager) {
	this.dispatch = dispatch
	this.client = NewFastHTTPClientOption(DefaultHTTPOption())
	this.intercepters = append(this.intercepters, &AccessLogIntercepter{})
	this.intercepters = append(this.intercepters, &LimitIntercepter{})
}

func (this *HttpProxy) AddIntercepter(i Intercepter) {
	if i != nil {
		this.intercepters = append(this.intercepters, i)
	}
}

func (this *HttpProxy) Start() {
	this.server = &fasthttp.Server{Handler: this.ServeHTTP}
	addr := fmt.Sprintf("0.0.0.0:%d", 80)
	this.server.ListenAndServe(addr)

}

func (this *HttpProxy) ServeHTTP(ctx *fasthttp.RequestCtx) {
	//获取访问的url
	url := string(ctx.Request.URI().RequestURI())
	domain := string(ctx.Request.Header.Host())
	//通过dispatch，获取api的定义
	api := this.dispatch.GetApi(domain, url)

	if api == nil { //不存在的时候的处理
		ctx.Response.SetBodyString("the url not found!")
		return
	}
	//根据api定义内容，进行 auth、access、等等的处理
	var res *fasthttp.Response
	if this.beforeCall(api, ctx) {
		//如果 cluster为"-",表示静态资源
		if api.Cluster.ID == Cluster_Static {
			res = this.getStaticResource(&ctx.Request, api.ServerUrl)
		} else {
			//根据分流和loadbalance选取server
			server := this.loadBalance(api, ctx)
			if server == nil {
				ctx.Response.SetBodyString("the server not exist!")
				return
			}

			//目标server调用
			res = this.call(ctx.Request, server, api.ServerUrl)

		}
	}

	//完成拦截器处理,处理服务器的返回值
	this.afterCall(api, ctx.ID(), res)

	//返回数据,回写response
	this.copyResponse(res, &ctx.Response)

	//defer,release调用
	defer this.releaseCall(api, ctx)

}

func (this *HttpProxy) beforeCall(api *meta.Api, ctx *fasthttp.RequestCtx) bool {
	for _, intercepter := range this.intercepters {
		ok, err := intercepter.Before(api, ctx)
		if !ok {
			err.Error()
			return false
		}
	}

	return true
}

func (this *HttpProxy) afterCall(api *meta.Api, id uint64, res *fasthttp.Response) {
	for _, intercepter := range this.intercepters {
		ok, err := intercepter.After(api, id, res)
		if !ok {
			err.Error()
			//TODO 错误处理

		}
	}

}

func (this *HttpProxy) releaseCall(api *meta.Api, ctx *fasthttp.RequestCtx) {

	for _, intercepter := range this.intercepters {
		ok, err := intercepter.Release(api, ctx)
		if !ok {
			err.Error()
		}
	}
}

//负载均衡
func (this *HttpProxy) loadBalance(api *meta.Api, ctx *fasthttp.RequestCtx) *meta.Server {

	if api != nil {
		context := GetRuntimeContext(api)
		//if context.QPS.Incr() {
		log.Println(context)
		if context.Lb != nil {
			log.Print(context.Lb)
			return context.Lb.Select(ctx, &api.Cluster.Servers)
		}
		//没有负载均衡设置，走random
		lservers := len(api.Cluster.Servers)
		if lservers > 1 {
			rand.Seed(time.Now().UnixNano())
			index := rand.Intn(lservers)
			return api.Cluster.Servers[index]
		} else if lservers > 0 {
			return api.Cluster.Servers[0]
		}
		//}

	}

	return nil
}

func (this *HttpProxy) call(req fasthttp.Request, server *meta.Server, url string) *fasthttp.Response {
	//需要进入重试处理
	r := copyRequest(&req)
	r.SetRequestURI("/" + url)
	res, err := this.client.Do(r, server.Addr(), nil)
	if err != nil {
		fmt.Println(err.Error())
		r := fasthttp.Response{}
		r.SetBodyString("the server error! used default return!")
		return &r
	}
	defer fasthttp.ReleaseRequest(r)
	return res
}

func (this *HttpProxy) copyResponse(source *fasthttp.Response, target *fasthttp.Response) {
	if source != nil {
		source.CopyTo(target)
		defer fasthttp.ReleaseResponse(source)
	}

}

func copyRequest(req *fasthttp.Request) *fasthttp.Request {
	newreq := fasthttp.AcquireRequest()
	newreq.Reset()
	req.CopyTo(newreq)
	return newreq
}

//处理静态资源
func (this *HttpProxy) getStaticResource(req *fasthttp.Request, url string) *fasthttp.Response {
	//TODO 完成从指定的静态目录中加载对应的url

	return nil
}
