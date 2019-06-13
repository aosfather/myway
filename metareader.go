package main

import (
	"fmt"
	"github.com/aosfather/myway/meta"
	"github.com/aosfather/myway/runtime"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

/**
  API及Server配置文件定义导入

*/

/**
  一个cluster 一个文件的导入
*/
type Cluster struct {
	Id            string   `yaml:"id"`
	Name          string   `yaml:"name"`
	Balance       int32    `yaml:"balance"`
	BalanceConfig string   `yaml:"balance_config"`
	Servers       []Server `yaml:"servers"`
}

func (this *Cluster) toServerCluster() *meta.ServerCluster {
	cluster := meta.ServerCluster{}
	cluster.ID = this.Id
	cluster.Name = this.Name
	cluster.Balance = meta.LoadBalance(this.Balance)
	cluster.BalanceConfig = this.BalanceConfig

	//加载server的定义
	for _, s := range this.Servers {
		cluster.Servers = append(cluster.Servers, s.toServer())
	}

	return &cluster
}

func (this *Cluster) Load(file string) {
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}

	err = yaml.Unmarshal(configFile, this)
	if err != nil {
		log.Fatalf("config file read error %v", err)
	}
}

type Server struct {
	Id   int64  `yaml:"id"`
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`
	Tag  string `yaml:"tag"`
}

func (this *Server) toServer() *meta.Server {
	s := meta.Server{}
	s.ID = this.Id
	s.Tag.Init(this.Tag)
	s.Ip = this.Ip
	s.Port = this.Port
	return &s
}

type apiFile struct {
	Namespace string      `yaml:"namespace"` //命名头
	Apis      []apiDefine `yaml:"apis"`      //api的定义
}

func (this *apiFile) Load(file string) {
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}

	err = yaml.Unmarshal(configFile, this)
	if err != nil {
		log.Fatalf("config file read error %v", err)
	}
}

type apiDefine struct {
	Url       string `yaml:"url"`
	Namespace string `yaml:"ns"`
	MaxQPS    int64  `yaml:"max_qps"`
	ServerUrl string `yaml:"server_url"`
	Cluster   string `yaml:"cluster"`
	Auth      bool   `yaml:"auth"`
}

func (this *apiDefine) toApi(fix string) *meta.Api {
	api := meta.Api{}
	api.ServerUrl = this.ServerUrl
	api.NameSpace = this.Namespace
	api.Url = "/" + fix + "/" + this.Url
	api.MaxQPS = this.MaxQPS
	if this.Auth {
		api.AuthFilter = "access_token"
	}
	return &api
}

//从文件中加载cluster的定义
func LoadClusterFromFile(path string, dispatch *runtime.DispatchManager) bool {
	files, e := ioutil.ReadDir(path)
	if e != nil {
		//TODO 输出错误信息，结束
		return false
	}

	for _, f := range files {
		fmt.Println(f.Name())
		filename := path + "/" + f.Name()
		c := Cluster{}
		c.Load(filename)
		fmt.Println(c)
		//注册集群
		dispatch.AddCluster(c.toServerCluster())
	}

	return true

}

//从文件中加载api的定义
func LoadAPIFromFile(path string, dispatch *runtime.DispatchManager) bool {
	files, e := ioutil.ReadDir(path)
	if e != nil {
		//TODO 输出错误信息，结束
		return false
	}

	for _, f := range files {
		fmt.Println(f.Name())
		filename := path + "/" + f.Name()
		apis := apiFile{}
		apis.Load(filename)
		fmt.Println(apis)
		//批量注册api
		for _, apidefine := range apis.Apis {
			api := apidefine.toApi(apis.Namespace)
			api.Cluster = dispatch.GetCluster(apidefine.Cluster)
			fmt.Println(api)
			fmt.Println(api.Cluster)
			dispatch.AddApi("", "", api)
		}

	}

	return true

}
