package runtime

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type PluginType byte

func ParsePluginType(pt string) PluginType {
	if pt == "balance" {
		return PT_Balance
	} else if pt == "" {

	}

	return PT_Balance
}

const (
	PT_Balance PluginType = 10 //负载均衡插件类型
	PT_Handle  PluginType = 11 //内部服务插件

)

//内部处理插件
type HandlePlugin func(req *fasthttp.Request) *fasthttp.Response

type pluginManager struct {
	plugins map[string]HandlePlugin
}

func (this *pluginManager) addPlugin(name string, plugin HandlePlugin) {
	if this.plugins == nil {
		this.plugins = make(map[string]HandlePlugin)
	}

	if name != "" && plugin != nil {
		this.plugins[name] = plugin
	}
}

func (this *pluginManager) callPlugin(name string, req *fasthttp.Request) *fasthttp.Response {
	fmt.Println("call plugin " + name)
	if name != "" && req != nil {
		p := this.plugins[name]
		if p != nil {
			return p(req)
		}

	}

	//构建通用错误
	res := fasthttp.Response{}
	res.SetBodyString("The plugin server not exist!")
	return &res
}
