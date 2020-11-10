package runtime

import (
	"github.com/valyala/fasthttp"
)

//request包装reader writer内容
type RequestWrap struct {
	Request *fasthttp.Request
}

//获取头
func (this *RequestWrap) GetHeader(key string) string {
	return string(this.Request.Header.Peek(key))
}

//获取参数
func (this *RequestWrap) GetParamter(key string) string {
	if string(this.Request.Header.Method()) == "GET" {
		return string(this.Request.URI().QueryArgs().Peek(key))
	}
	return string(this.Request.PostArgs().Peek(key))

}

//获取内容
func (this *RequestWrap) GetBody() []byte {
	return this.Request.Body()
}

//设置头
func (this *RequestWrap) SetHeader(key, value string) {
	this.Request.Header.Set(key, value)
}

//删除header 头
func (this *RequestWrap) RemoveHeader(key string) {
	this.Request.Header.Del(key)
}

//设置参数
func (this *RequestWrap) SetParamter(key string, value string) {
	var args *fasthttp.Args
	if string(this.Request.Header.Method()) == "GET" {
		args = this.Request.URI().QueryArgs()
	} else {
		args = this.Request.PostArgs()
	}

	if args.Has(key) {
		args.Set(key, value)
	} else {
		args.Add(key, value)
	}
}

//删除参数
func (this *RequestWrap) RemoveParamter(key string) {
	var args *fasthttp.Args
	if string(this.Request.Header.Method()) == "GET" {
		args = this.Request.URI().QueryArgs()
	} else {
		args = this.Request.PostArgs()
	}
	args.Del(key)
}

//设置body体
func (this *RequestWrap) SetBody(b []byte) {
	this.Request.SetBody(b)
}

//response包装reader writer
type ResponseWrap struct {
	Response *fasthttp.Response
}

//获取头
func (this *ResponseWrap) GetHeader(key string) string {
	return string(this.Response.Header.Peek(key))
}

//获取参数,不支持这个
func (this *ResponseWrap) GetParamter(key string) string {
	return ""
}

//获取内容
func (this *ResponseWrap) GetBody() []byte {
	return this.Response.Body()
}

//设置头
func (this *ResponseWrap) SetHeader(key, value string) {
	this.Response.Header.Set(key, value)
}

//删除header 头
func (this *ResponseWrap) RemoveHeader(key string) {
	this.Response.Header.Del(key)
}

//设置参数，不支持
func (this *ResponseWrap) SetParamter(key string, value string) {

}

//删除参数，不支持
func (this *ResponseWrap) RemoveParamter(key string) {

}

//设置body体
func (this *ResponseWrap) SetBody(b []byte) {
	this.Response.SetBody(b)
}
