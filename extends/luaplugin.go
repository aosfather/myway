package extends

import (
	"encoding/json"
	"fmt"
	"github.com/aosfather/myway/meta"
	"github.com/aosfather/myway/runtime"
	"github.com/valyala/fasthttp"
	"github.com/yuin/gopher-lua"
	"sync"
)

/**
  插件
  基于lua脚本的插件
  1、内部的方法
  2、功能的扩展
*/
type PluginCore interface {
	GetFromContext(key string) interface{} //获取信息
	SetContext(key string, v interface{})  //设置信息
}
type LuaPlugin struct {
	l          *lua.LState
	code       string //脚本名（后缀.lua）
	catalog    string //分类目录
	locker     sync.RWMutex
	pluginType runtime.PluginType //插件类型
	core       PluginCore         //插件核心
}

func (this *LuaPlugin) Init(catalog, code, ptype string) {
	this.catalog = catalog
	this.code = code
	this.pluginType = runtime.ParsePluginType(ptype)
	this.l = lua.NewState()
	//加载代码
	filename := catalog + "/" + code + ".lua"
	this.l.DoFile(filename)
	//函数注册

}

func (this *LuaPlugin) GetType() runtime.PluginType {
	return this.pluginType
}

//注册函数给脚本使用
func (this *LuaPlugin) registerFunction(name string, function lua.LGFunction) {
	this.l.SetGlobal(name, this.l.NewFunction(function))

}

//运行脚本
func (this *LuaPlugin) run(method string, p interface{}) string {
	this.locker.Lock()
	fn := this.l.GetGlobal(method)
	err := this.l.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, this.toTable(p))

	if err != nil {
		fmt.Println("call apply")
		fmt.Println(err)
		return err.Error()
	}

	title := this.l.Get(-1).String()
	fmt.Println(title)
	this.l.Pop(1)

	defer this.locker.Unlock()
	return ""
}

func (this *LuaPlugin) toTable(event interface{}) *lua.LTable {
	if value, ok := event.(string); ok == true {
		return this.jsonToTable(value)
	} else if value, ok := event.(map[string]string); ok == true {
		return this.mapToTable(value)
	} else {
		t := this.l.NewTable()
		return t
	}
}

func (this *LuaPlugin) jsonToTable(str string) *lua.LTable {
	data := make(map[string]string)
	json.Unmarshal([]byte(str), &data)
	return this.mapToTable(data)
}

func (this *LuaPlugin) mapToTable(data map[string]string) *lua.LTable {
	t := this.l.NewTable()
	for key, value := range data {
		fmt.Println(key, value)
		t.RawSet(lua.LString(key), lua.LString(value))
	}
	return t
}

//
type BalancePlugin struct {
	LuaPlugin
}

func (this *BalancePlugin) Config(p string) {
	this.run("Config", p)

}

func (this *BalancePlugin) Select(req *fasthttp.RequestCtx, servers *[]*meta.Server) *meta.Server {

	return nil
}
