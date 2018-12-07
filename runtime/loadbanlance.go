package runtime

import (
	"github.com/aosfather/myway/meta"
	"github.com/valyala/fasthttp"
	"hash/crc32"
	"strings"
	"sync/atomic"
)

/*
  负载均衡器
*/
type LoadBalance interface {
	Config(p string)
	Select(req *fasthttp.RequestCtx, servers *[]*meta.Server) *meta.Server
}

type abstractLoadBalance struct {
}

func (this *abstractLoadBalance) Config(p string) {

}

/**
  随机均衡负载器
*/
type RoundRobin struct {
	abstractLoadBalance
	ops *uint64
}

func NewRoundRobin() LoadBalance {
	var ops uint64 = 0
	return &RoundRobin{
		ops: &ops,
	}
}

// Select select a server from servers using RoundRobin
func (rr RoundRobin) Select(req *fasthttp.RequestCtx, servers *[]*meta.Server) *meta.Server {
	l := uint64(len(*servers))
	if 0 >= l {
		return nil
	}
	index := int(atomic.AddUint64(rr.ops, 1) % l)
	return (*servers)[index]
}

/**
  ip hash loadbalance
*/
type IPHash struct {
	abstractLoadBalance
}

func (rr IPHash) Select(req *fasthttp.RequestCtx, servers *[]*meta.Server) *meta.Server {
	l := len(*servers)
	//取ip地址的hash
	ip := GetRealClientIP(req)
	hc := Hashcode(ip)
	index := int(hc % l)
	return (*servers)[index]
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

func (this *TagLoadBalance) Config(p string) {
	this.Tag = p
	this.balance = buildBalance(0)
}
func (this *TagLoadBalance) Select(req *fasthttp.RequestCtx, servers *[]*meta.Server) *meta.Server {
	var tagServers []*meta.Server
	//获取有指定tag标签的服务器列表
	for _, v := range *servers {
		if v.Tag.Has(this.Tag) {
			tagServers = append(tagServers, v)
		}

	}

	//然后使用iphash方式，进行分配
	l := len(tagServers)
	if l == 0 {
		return nil
	}

	//取ip地址的hash
	ip := GetRealClientIP(req)
	hc := Hashcode(ip)
	index := int(hc % l)

	return tagServers[index]
}

//工厂方法，构建负载均衡器
func buildBalance(b meta.LoadBalance) LoadBalance {
	switch b {
	case meta.LBIPHash:
		return new(IPHash)
	case meta.LBRoundRobin:
		return NewRoundRobin()
	case meta.LBTag:
		return new(TagLoadBalance)
	default:
		return new(IPHash)

	}
}
