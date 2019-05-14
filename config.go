package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

/*
 配置文件模型--YAML
 1、server 的配置文件夹
 2、api 的配置文件夹
 3、
*/

type yamlConfig struct {
	Config ApplicationConfigurate `yaml:"path"`
}
type ApplicationConfigurate struct {
	ServerPath string `yaml:"server"`
	ApiPath    string `yaml:"api"`
}

//从文件装载信息
func (this *yamlConfig) Load(file string) {
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}

	err = yaml.Unmarshal(configFile, this)
	if err != nil {
		log.Fatalf("config file read error %v", err)
	}
}
