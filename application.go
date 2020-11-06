package main

import (
	"fmt"
	"github.com/aosfather/myway/meta"
	"github.com/aosfather/myway/runtime"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

/**
  应用模型
*/
type Application struct {
	Version string //版本
	Module  string //虚拟模块
	Host    string //对应主机
	Prefix  string
	Apis    []ApplicationApi //接口列表
}

func (this *Application) ToApplicationMapper() *meta.ApplicationMapper {
	fmt.Println(this)
	mapper := meta.ApplicationMapper{App: this.Module, Host: this.Host}
	//解析出主机地址和端口
	host, port := extractHostAndPort(this.Host)
	cluster := meta.ServerCluster{}
	cluster.AddServer(&meta.Server{ID: 1, Ip: host, Port: port, MaxQPS: 1000})
	mapper.Cluster = &cluster
	for _, apidefine := range this.Apis {
		mapper.AddMapper(apidefine.ToApi(this))
	}

	return &mapper
}

func extractHostAndPort(url string) (string, int) {
	var host string
	var port int
	//解析host
	host = strings.ReplaceAll(url, "http://", "")
	uris := strings.Split(host, "/")
	if len(uris) >= 1 {
		host = uris[0]
	}

	//获取port信息
	tmp := strings.Split(host, ":")
	if len(tmp) == 1 {
		port = 80
	} else {
		host = tmp[0]
		port, _ = strconv.Atoi(tmp[1])
	}
	return host, port
}

func (this *Application) Load(file string) {
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}

	err = yaml.Unmarshal(configFile, this)
	if err != nil {
		log.Fatalf("config file read error %v", err)
	}
}

//应用接口
type ApplicationApi struct {
	Url       string //api的url
	TargetUrl string `yaml:"targetUrl"` //目标地址
	Label     string //接口标题
	Desc      string //描述
}

func (this *ApplicationApi) ToApi(app *Application) *meta.ApiMapper {
	api := meta.ApiMapper{}
	api.Url = "/" + app.Module + "/" + this.Url
	api.TargetUrl = app.Prefix + this.TargetUrl
	fmt.Println(this)
	fmt.Println(api)
	return &api
}

func LoadApplicationFromFile(path string, dispatch *runtime.DispatchManager) bool {
	files, e := ioutil.ReadDir(path)
	if e != nil {
		//TODO 输出错误信息，结束
		return false
	}

	for _, f := range files {
		fmt.Println(f.Name())
		filename := path + "/" + f.Name()
		apis := Application{}
		apis.Load(filename)
		dispatch.AddApplication(apis.ToApplicationMapper())
	}

	return true

}
