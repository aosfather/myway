package main

import (
	"github.com/aosfather/bingo_utils/codes"
	"github.com/aosfather/myway/filter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"regexp"
)

type ContextLoad func(file string)

var loadMap = make(map[string]ContextLoad)

//从yaml文件中加载数据
func loadFromYaml(file string, target interface{}) {
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}

	err = yaml.Unmarshal(configFile, target)
	if err != nil {
		log.Fatalf("config file read error %v", err)
	}
}

//脱敏相关的定义
type desensitationConfigFile struct {
	Name     string
	Version  string
	Contexts []desensitationConfig
}

func (this *desensitationConfigFile) load(file string) {
	loadFromYaml(file, this)
}

type desensitationConfig struct {
	Name  string
	Label string
	Data  desensitationDef
}

type desensitationDef struct {
	Head    int
	Tail    int
	Padding bool
	Replace bool
	Pattern string
	Mask    rune
}

func (this *desensitationDef) create() *codes.Sensitive {
	var p *regexp.Regexp = nil
	if this.Pattern != "" {
		p = regexp.MustCompile(this.Pattern)
	}
	return &codes.Sensitive{this.Head, this.Tail, this.Padding, this.Replace, p, this.Mask}
}

//加载脱敏的定义
func loadDesensitation(file string) {
	cf := desensitationConfigFile{}
	cf.load(file)
	//注册脱敏处理器
	for _, cv := range cf.Contexts {
		filter.RegisterDesensitive(cv.Name, cv.Data.create())
	}

}

//注册context loader
func init() {
	loadMap["desensitation"] = loadDesensitation
}
