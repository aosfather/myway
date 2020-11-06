package runtime

import (
	"fmt"
	"github.com/aosfather/myway/meta"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

//运行上下文
type RunTimeContext map[string]interface{}

/**
  拦截器接口
*/

type Intercepter interface {
	Before(api *meta.ApiMapper, ctx *fasthttp.RequestCtx, context RunTimeContext) (bool, error) //调用之前
	After(api *meta.ApiMapper, ctx *fasthttp.Response, context RunTimeContext) (bool, error)    //调用完成之后
	Release(api *meta.ApiMapper, ctx *fasthttp.RequestCtx) (bool, error)                        //网关完成返回处理完之后释放资源
}

//抽象拦截器
type BaseIntercepter struct {
}

func (this *BaseIntercepter) Before(api *meta.ApiMapper, ctx *fasthttp.RequestCtx, context RunTimeContext) (bool, error) {
	return true, nil
}

func (this *BaseIntercepter) After(api *meta.ApiMapper, ctx *fasthttp.Response, context RunTimeContext) (bool, error) {
	return true, nil
}

func (this *BaseIntercepter) Release(api *meta.ApiMapper, ctx *fasthttp.RequestCtx) (bool, error) {
	return true, nil
}

//鉴权拦截器
type AuthIntercepter struct {
	BaseIntercepter
}

func (this *AuthIntercepter) Before(api *meta.ApiMapper, ctx *fasthttp.RequestCtx, context RunTimeContext) (bool, error) {
	return true, nil
}

//限流拦截器
type LimitIntercepter struct {
	BaseIntercepter
}

func (this *LimitIntercepter) Before(api *meta.ApiMapper, ctx *fasthttp.RequestCtx, context RunTimeContext) (bool, error) {
	c := GetRuntimeValve(api)
	if c.QPS.Incr() {
		return true, nil
	}
	ctx.Response.SetBodyString("the api limit max qps!")
	return false, fmt.Errorf("over the max qps")
}

////访问拦截器
//type AccessIntercepter struct {
//	BaseIntercepter
//	tm core.TokenManager
//	rm core.RoleManager
//}
//
//func (this *AccessIntercepter) Init(addr string, db int, expire int64, rmm core.RoleMetaManager) {
//	ts := RedisTokenStore{}
//	ts.Init(addr, db)
//	this.tm.Store = &ts
//	if expire <= 0 {
//		this.tm.Expire = 7200 //默认两小时
//	} else {
//		this.tm.Expire = expire
//	}
//
//	//设置读取权限元信息
//	this.rm.SetMetaManager(rmm)
//}
//
//func (this *AccessIntercepter) Before(api *meta.Api, ctx *fasthttp.RequestCtx) (bool, error) {
//	//1、看api是否需要token
//	if api.AuthFilter == "access_token" {
//		fmt.Println("in this")
//		//2、看token是否有效
//		token := string(ctx.Request.PostArgs().Peek("access_token"))
//		fmt.Println("in this token")
//		t := this.tm.GetToken(token)
//		//3、看是否有权限调用
//		if t != nil {
//			m := make(map[string]string)
//			m["url"] = api.Url
//			if this.rm.Validate(t.Role, "url", m) {
//				return true, nil
//			}
//
//		}
//
//		ctx.Response.SetBodyString("no auth!")
//		return false, fmt.Errorf("no auth accssee the url!")
//
//	}
//
//	return true, nil
//}
//
//type RedisTokenStore struct {
//	client *redis.Client
//}
//
//func (this *RedisTokenStore) Init(addr string, db int) {
//	this.client = redis.NewClient(&redis.Options{
//		Addr:     addr,
//		Password: "", // no password set
//		DB:       db,
//	})
//}
//
////获取token
//func (this *RedisTokenStore) GetToken(id string) *core.AccessToken {
//	if this.client != nil {
//
//	}
//	return nil
//}
//
////保存token
//func (this *RedisTokenStore) SaveToken(id string, t *core.AccessToken, expire int64) {
//	if this.client != nil {
//
//	}
//}

//访问日志
type AccessLogIntercepter struct {
	BaseIntercepter
}

func (this *AccessLogIntercepter) Before(api *meta.ApiMapper, ctx *fasthttp.RequestCtx, context RunTimeContext) (bool, error) {
	log.Println("call" + api.Url + fmt.Sprintf("%d", ctx.ID()))
	context["__begin__"] = time.Now().UnixNano()
	return true, nil
}

func (this *AccessLogIntercepter) After(api *meta.ApiMapper, ctx *fasthttp.Response, context RunTimeContext) (bool, error) {
	thenow := time.Now().UnixNano()
	context["__end__"] = thenow
	start := context["__begin__"].(int64)
	used := thenow - start
	context["__used__"] = used
	log.Println("after call " + api.Url + fmt.Sprintf(" used:%d ms", (time.Duration(used)/time.Millisecond)))
	return true, nil
}
