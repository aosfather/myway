package runtime

import (
	"container/list"
	"github.com/valyala/fasthttp"
	"hash/crc32"
	"strings"
	"sync/atomic"
)

/*
  负载均衡器
*/
type LoadBalance interface {
	Select(req *fasthttp.RequestCtx, servers *list.List) int
}

/**
  随机均衡负载器
*/
type RoundRobin struct {
	ops *uint64
}

func NewRoundRobin() LoadBalance {
	var ops uint64 = 0
	return RoundRobin{
		ops: &ops,
	}
}

// Select select a server from servers using RoundRobin
func (rr RoundRobin) Select(req *fasthttp.RequestCtx, servers *list.List) int {
	l := uint64(servers.Len())
	if 0 >= l {
		return -1
	}
	return int(atomic.AddUint64(rr.ops, 1) % l)
}

/**
  ip hash loadbalance
*/
type IPHash struct {
}

func (rr IPHash) Select(req *fasthttp.RequestCtx, servers *list.List) int {
	l := servers.Len()
	//取ip地址的hash
	ip := GetRealClientIP(req)
	hc := Hashcode(ip)
	return int(hc % l)
}

func GetRealClientIP(ctx *fasthttp.RequestCtx) string {
	xforward := ctx.Request.Header.Peek("X-Forwarded-For")
	if nil == xforward {
		return strings.SplitN(ctx.RemoteAddr().String(), ":", 2)[0]
	}

	return strings.SplitN(string(xforward), ",", 2)[0]
}

func Hashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

//根据标签进行处理
type TagLoadBalance struct {
	Tag     string      //tag
	balance LoadBalance //使用的真实负载策略
}

func (this *TagLoadBalance) Select(req *fasthttp.RequestCtx, servers *list.List) int {
	//获取有指定tag标签的服务器列表

	//然后使用iphash方式，进行分配
	l := servers.Len()
	//取ip地址的hash
	ip := GetRealClientIP(req)
	hc := Hashcode(ip)
	return int(hc % l)
}
