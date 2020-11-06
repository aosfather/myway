package runtime

import (
	"fmt"
	"github.com/aosfather/myway/meta"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"strings"
	"time"
)

const Cluster_Static = "-" //静态资源
const Cluster_This = "."   //本地插件

type HttpConfig struct {
	Root string
	Port int
}

//基本的代理
type HttpProxy struct {
	port         int
	server       *fasthttp.Server
	dispatch     *DispatchManager //分发管理
	intercepters []Intercepter    //拦截器处理
	client       *FastHTTPClient
	plugins      *pluginManager
	staticRoot   string //静态资源根目录
}

func (this *HttpProxy) Init(dispatch *DispatchManager) {
	this.dispatch = dispatch
	this.client = NewFastHTTPClientOption(DefaultHTTPOption())
	this.intercepters = append(this.intercepters, &AccessLogIntercepter{})
	this.intercepters = append(this.intercepters, &LimitIntercepter{})
	this.plugins = &pluginManager{}
}

func (this *HttpProxy) AddPlugin(name string, plugin HandlePlugin) {
	if this.plugins != nil {
		this.plugins.addPlugin(name, plugin)
	}
}

func (this *HttpProxy) AddIntercepter(i Intercepter) {
	if i != nil {
		this.intercepters = append(this.intercepters, i)
	}
}

//配置
func (this *HttpProxy) SetConfig(config HttpConfig) {
	this.staticRoot = config.Root
	this.port = config.Port
}

func (this *HttpProxy) Start() {
	this.server = &fasthttp.Server{Handler: this.ServeHTTP}

	if this.port <= 0 {
		this.port = 80
	}

	addr := fmt.Sprintf("0.0.0.0:%d", this.port)
	this.server.ListenAndServe(addr)

}

/**
  运作流程
   1、获取绑定的url对应的API
   2、如果没有，查找映射的app（注册中心模式）
   3、
*/
func (this *HttpProxy) ServeHTTP(ctx *fasthttp.RequestCtx) {
	//获取访问的url
	url := string(ctx.Request.URI().RequestURI())
	domain := string(ctx.Request.Header.Host())
	//通过dispatch，获取api的定义
	fmt.Println(domain, "=>", url)
	api := this.getApi(domain, url) //this.dispatch.GetApi(domain, url)
	if api == nil {
		ctx.Response.SetBodyString("the url not found!")
		return
	}

	//根据api定义内容，进行 auth、access、等等的处理
	var res *fasthttp.Response
	context := make(RunTimeContext)

	if this.beforeCall(api, ctx, context) {
		//根据分流和loadbalance选取server
		server := this.loadBalance(api, ctx)
		if server == nil {
			ctx.Response.SetBodyString("the server not exist!")
			return
		} else {
			req := ctx.Request
			//目标server调用
			nreq := this.inputFilter(&req, server.Addr(), api, context)
			res = this.proxyTarget(nreq)
		}

	}
	//完成拦截器处理,处理服务器的返回值
	this.afterCall(api, res, context)
	//返回数据,回写response
	this.outputFilter(res, &ctx.Response, api, context)

	//defer,release调用
	defer this.releaseCall(api, ctx)

}

func (this *HttpProxy) getApi(domain, url string) *meta.ApiMapper {
	api := this.dispatch.GetApi(domain, url)

	if api == nil { //不存在的时候的处理
		//看是否是应用映射
		appname := url[1:]
		end := strings.Index(appname, "/")
		var targetUri string
		if end > 0 {
			targetUri = appname[end:]
			appname = appname[0:end]
		}

		appmapper := this.dispatch.GetApplication(appname)
		if appmapper != nil {
			api = &meta.ApiMapper{Url: url, TargetUrl: targetUri}
			appmapper.AddMapper(api)
			this.dispatch.AddApi("", "", api)
		}
	}

	return api
}

func (this *HttpProxy) beforeCall(api *meta.ApiMapper, ctx *fasthttp.RequestCtx, context RunTimeContext) bool {

	for _, intercepter := range this.intercepters {
		ok, err := intercepter.Before(api, ctx, context)
		if !ok {
			err.Error()
			return false
		}
	}

	return true
}

func (this *HttpProxy) afterCall(api *meta.ApiMapper, res *fasthttp.Response, context RunTimeContext) {
	for _, intercepter := range this.intercepters {
		ok, err := intercepter.After(api, res, context)
		if !ok {
			err.Error()
			//TODO 错误处理

		}
	}

}

func (this *HttpProxy) releaseCall(api *meta.ApiMapper, ctx *fasthttp.RequestCtx) {

	for _, intercepter := range this.intercepters {
		ok, err := intercepter.Release(api, ctx)
		if !ok {
			err.Error()
		}
	}
}

//负载均衡
func (this *HttpProxy) loadBalance(api *meta.ApiMapper, ctx *fasthttp.RequestCtx) *meta.Server {

	if api != nil {
		context := GetRuntimeValve(api)
		//if context.QPS.Incr() {
		log.Println(context)
		if context.Lb != nil {
			log.Print(context.Lb)
			return context.Lb.Select(ctx, &api.GetCluster().Servers)
		}
		//没有负载均衡设置，走random
		if api.GetCluster() == nil {
			return nil
		}
		lservers := len(api.GetCluster().Servers)
		if lservers > 1 {
			rand.Seed(time.Now().UnixNano())
			index := rand.Intn(lservers)
			return api.GetCluster().Servers[index]
		} else if lservers > 0 {
			return api.GetCluster().Servers[0]
		}
		//}

	}

	return nil
}

func (this *HttpProxy) inputFilter(req *fasthttp.Request, addr string, api *meta.ApiMapper, context RunTimeContext) *fasthttp.Request {
	newreq := fasthttp.AcquireRequest()
	newreq.Reset()
	req.CopyTo(newreq)
	//需要进入重试处理
	newreq.SetRequestURI("/" + api.TargetUrl)
	newreq.SetHost(addr)
	//进行filter处理

	return newreq
}

//调用目标服务
func (this *HttpProxy) proxyTarget(req *fasthttp.Request) *fasthttp.Response {
	res, err := this.client.Do(req, string(req.Host()), nil)
	if err != nil {
		fmt.Println(err.Error())
		r := fasthttp.Response{}
		r.SetBodyString("the server error! used default return!")
		return &r
	}

	defer fasthttp.ReleaseRequest(req)
	return res
}

//输出结果处理
func (this *HttpProxy) outputFilter(source, target *fasthttp.Response, api *meta.ApiMapper, context RunTimeContext) {
	if source != nil {
		source.Header.Del("Transfer-Encoding")
		source.CopyTo(target)
		defer fasthttp.ReleaseResponse(source)
	}
}
