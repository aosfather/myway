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
	System SystemConfigurate      `yaml:"system"`
}
type ApplicationConfigurate struct {
	ServerPath    string   `yaml:"server"`
	AppPath       string   `yaml:"app"`
	UserPath      string   `yaml:"user"`
	RolePath      string   `yaml:"role"`
	AccessLogFile string   `yaml:"access_log"`
	logFile       string   `yaml:"log"`
	Eureka        []string `yaml:"eureka"`
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

type SystemConfigurate struct {
	AuthRedis AuthRedisConfig `yaml:"auth_redis"`
}

type AuthRedisConfig struct {
	Address  string `yaml:"address"`
	DataBase int    `yaml:"db"`
	Expire   int64  `yaml:"expire"`
}
