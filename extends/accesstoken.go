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
	fmt.Println(req.PostArgs())
	fmt.Println(string(req.RequestURI()))
	res := fasthttp.Response{}
	//读取参数
	var grandType, user, secret string
	res.SetBodyString("hello!")

	if len(this.builder) > 0 {
		for _, tb := range this.builder {
			if tb.GetGrandType() == grandType {
				b, token := tb.BuildToken(user, secret)
				if b {
					fmt.Println(token)
					res.SetBodyString(fmt.Sprintf(`{access_token:"%s"}`, token))
				}
			}
		}
	}

	return &res
}
