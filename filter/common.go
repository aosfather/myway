package filter

import (
	"log"
)

/**
  访问管道filter处理
   这个概念借的是logstash的概念，特指数据处理这块内容
   包括：
    参数修改
    modify[参数修改]： 参数改名、删除参数、参数移动、参数内容修改
    add[参数新增]：环境变量里的内容：例如用户名等、当前日期、当前时间、header内容
     header修改
      modify：移除header、修改header内容
      add header
*/

//操作实体对象
type EntityReader interface {
	GetHeader(key string) string   //获取头
	GetParamter(key string) string //获取参数
	GetBody() []byte               //获取内容
}

//实体写
type EntityWriter interface {
	SetHeader(key, value string)          //设置头
	RemoveHeader(key string)              //删除header 头
	SetParamter(key string, value string) //设置参数
	RemoveParamter(key string)            //删除参数
	SetBody(b []byte)                     //设置body体
}

//filter定义
type Filter interface {
	Name() string
	DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error
}

//过滤器链
type FilterChain []Filter

func (this *FilterChain) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) {
	chain := *this
	if len(chain) > 0 {
		for index, f := range chain {
			log.Println(index)
			if f != nil {
				f.DoFilter(r, w, context)
			} else {
				log.Println("the filter is nil!")
			}

		}
	}
}

//filter管理器
type FilterManager struct {
	filters map[string]Filter
}

func (this *FilterManager) Init() {
	this.filters = make(map[string]Filter)
}

func (this *FilterManager) Register(id string, f Filter) {
	if f != nil {
		this.filters[id] = f
	}
}

func (this *FilterManager) CreateChain(ids ...string) FilterChain {
	f := FilterChain{}
	for _, n := range ids {
		f = append(f, this.filters[n])
	}

	return f
}
