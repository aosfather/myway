package extends

import (
	"fmt"
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
}

type AccessTokenImp struct {
	builder []TokenCreator
}

func (this *AccessTokenImp) Add(tb TokenCreator) {
	if tb != nil {
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
