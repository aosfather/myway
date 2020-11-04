package extends

import (
	"fmt"
	"github.com/aosfather/myway/core"
	"github.com/aosfather/myway/meta"
	"github.com/aosfather/myway/runtime"
	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
)

/**
  基本的Access toke 实现
  参数：grandtype ==access_token ,user==分配的用户或应用id,secret==秘钥
*/
//创建token的creator
type TokenCreator interface {
	GetGrandType() string
	BuildToken(user string, secret string) (bool, string)
	SetTokenManager(tm *core.TokenManager)
}

type AccessTokenImp struct {
	runtime.BaseIntercepter
	tm      core.TokenManager
	rm      core.RoleManager
	builder []TokenCreator
}

func (this *AccessTokenImp) Add(tb TokenCreator) {
	if tb != nil {
		tb.SetTokenManager(&this.tm)
		this.builder = append(this.builder, tb)
	}
}
func (this *AccessTokenImp) Call(req *fasthttp.Request) *fasthttp.Response {
	res := fasthttp.Response{}

	args := req.URI().QueryArgs()

	//读取参数
	var grandType, user, secret string
	grandType = string(args.Peek("grandtype"))
	user = string(args.Peek("user"))
	secret = string(args.Peek("secret"))
	if grandType == "" || user == "" {
		res.SetBodyString("the parameter grandtype and user is must!")
	} else {
		tb := this.getTokenCreator(grandType)
		if tb == nil {
			res.SetBodyString("the grandtype not surport!")
		} else {
			b, token := tb.BuildToken(user, secret)
			if b {
				fmt.Println(token)
				res.SetBodyString(fmt.Sprintf(`{access_token:"%s"}`, token))
			} else {
				res.SetBodyString("access deny!")
			}
		}
	}

	return &res
}

func (this *AccessTokenImp) getTokenCreator(t string) TokenCreator {
	if len(this.builder) > 0 {
		for _, tb := range this.builder {
			if tb.GetGrandType() == t {

				return tb
			}
		}
	}

	return nil
}

func (this *AccessTokenImp) Init(addr string, db int, expire int64, rmm core.RoleMetaManager) {
	ts := RedisTokenStore{}
	ts.Init(addr, db)
	this.tm.Store = &ts
	if expire <= 0 {
		this.tm.Expire = 7200 //默认两小时
	} else {
		this.tm.Expire = expire
	}

	//设置读取权限元信息
	this.rm.SetMetaManager(rmm)
}

func (this *AccessTokenImp) Before(api *meta.ApiMapper, ctx *fasthttp.RequestCtx) (bool, error) {
	//1、看api是否需要token
	//if api.AuthFilter == "access_token" {
	//	fmt.Println("in this")
	//	//2、看token是否有效
	//	token := string(ctx.Request.PostArgs().Peek("access_token"))
	//	fmt.Println("in this token")
	//	t := this.tm.GetToken(token)
	//	//3、看是否有权限调用
	//	if t != nil {
	//		m := make(map[string]string)
	//		m["url"] = api.Url
	//		if this.rm.Validate(t.Role, "url", m) {
	//			return true, nil
	//		}
	//
	//	}
	//
	//	ctx.Response.SetBodyString("no auth!")
	//	return false, fmt.Errorf("no auth accssee the url!")
	//
	//}

	return true, nil
}

type RedisTokenStore struct {
	client *redis.Client
}

func (this *RedisTokenStore) Init(addr string, db int) {
	this.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       db,
	})
}

//获取token
func (this *RedisTokenStore) GetToken(id string) *core.AccessToken {
	if this.client != nil {

	}
	return nil
}

//保存token
func (this *RedisTokenStore) SaveToken(id string, t *core.AccessToken, expire int64) {
	if this.client != nil {

	}
}
