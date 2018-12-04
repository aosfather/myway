package runtime

import (
	"github.com/aosfather/myway/meta"
	"github.com/valyala/fasthttp"
)

/**
  拦截器接口
*/

type Intercepter interface {
	Before(api *meta.Api, ctx *fasthttp.RequestCtx) (bool, error)  //调用之前
	After(api *meta.Api, ctx *fasthttp.Response) (bool, error)     //调用完成之后
	Release(api *meta.Api, ctx *fasthttp.RequestCtx) (bool, error) //网关完成返回处理完之后释放资源
}

//抽象拦截器
type baseIntercepter struct {
}

func (this *baseIntercepter) Before(api *meta.Api, ctx *fasthttp.RequestCtx) (bool, error) {
	return true, nil
}

func (this *baseIntercepter) After(api *meta.Api, ctx *fasthttp.Response) (bool, error) {
	return true, nil
}

func (this *baseIntercepter) Release(api *meta.Api, ctx *fasthttp.RequestCtx) (bool, error) {
	return true, nil
}

//鉴权拦截器
type AuthIntercepter struct {
	baseIntercepter
}

func (this *AuthIntercepter) Before(api *meta.Api, ctx *fasthttp.RequestCtx) (bool, error) {
	return true, nil
}

//限流拦截器
type LimitIntercepter struct {
}

func (this *LimitIntercepter) Before(api *meta.Api, ctx *fasthttp.RequestCtx) (bool, error) {
	return false, nil
}

//访问拦截器
type AccessIntercepter struct {
}

func (this *AccessIntercepter) Before(api *meta.Api, ctx *fasthttp.RequestCtx) (bool, error) {

	return false, nil
}
