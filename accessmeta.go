package main

import (
	"fmt"
	"github.com/aosfather/myway/core"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

//用户定义
type userInfo struct {
	Name     string `yaml:"name"`
	PassWord string `yaml:"pwd"`
	Role     string `yaml:"role"`
}

func (this *userInfo) Load(file string) {
	LoadFromYaml(this, file)
}

//角色定义
type roleInfo struct {
	Name      string    `yaml:"name"`
	SuperRole string    `yaml:"super"` //父级角色
	White     []urlInfo `yaml:"white"`
	Black     []urlInfo `yaml:"black"`
}

func (this *roleInfo) Load(file string) {
	LoadFromYaml(this, file)
}

//权限url
type urlInfo struct {
	Url string `yaml:"url"`
}

//权限元数据管理
type yamlAuthMetaManager struct {
	class *core.AuthClass
	roles map[string]*core.Role
}

func (this *yamlAuthMetaManager) Load(rolepath string) {
	files, e := ioutil.ReadDir(rolepath)
	if e != nil {
		//TODO 输出错误信息，结束
		return
	}

	if this.roles == nil {
		this.roles = make(map[string]*core.Role)
	}

	//权限class，只有url字段。名称为access
	if this.class == nil {
		this.class = core.CreateAuthClass("asscess")
		this.class.AddFieldByParame("url", "地址", core.FT_STRING)
	}

	for _, f := range files {
		fmt.Println(f.Name())
		filename := rolepath + "/" + f.Name()
		role := roleInfo{}
		role.Load(filename)
		fmt.Println(role)

		//将角色定义转换成标准的权限对象
		r := core.Role{}
		r.Code = role.Name
		//TODO 后续加入处理父类角色的功能
		aset := core.AuthobjectSet{}
		aset.Class = this.class
		r.AddAuthObjectSet(&aset)

		//白名单
		for _, u := range role.White {
			aset.AddAuthObject(false, []*core.AuthFieldValue{aset.BuildFieldValue("url", core.VT_SINGLE, u.Url)})
		}

		//黑名单
		for _, u := range role.Black {
			aset.AddAuthObject(true, []*core.AuthFieldValue{aset.BuildFieldValue("url", core.VT_SINGLE, u.Url)})
		}

		//加入角色池里
		this.roles[r.Code] = &r

	}
}

//用户管理
type userManager struct {
	users map[string]*userInfo
	tm    *core.TokenManager
}

func (this *userManager) Load(userpath string) {
	fmt.Println(userpath)
	files, e := ioutil.ReadDir(userpath)
	if e != nil {
		//TODO 输出错误信息，结束
		fmt.Println(e.Error())
		return
	}

	if this.users == nil {
		this.users = make(map[string]*userInfo)
	}

	for _, f := range files {
		fmt.Println(f.Name())
		filename := userpath + "/" + f.Name()
		user := userInfo{}
		user.Load(filename)
		fmt.Println(user)
		this.users[user.Name] = &user

	}
}

func (this *userManager) SetTokenManager(tm *core.TokenManager) {
	this.tm = tm
}

func (this *userManager) GetGrandType() string {
	return "access_token"
}
func (this *userManager) BuildToken(user string, secret string) (bool, string) {
	v, r := this.Validate(user, secret)
	if v {
		//调用tokenmanager创建token
		if this.tm != nil {
			_, id := this.tm.CreateToken(user, r)
			return true, id
		}
		fmt.Println("tm is nil")
	}

	return false, "no token"
}

func (this *userManager) Validate(name, pwd string) (bool, string) {
	if name == "" || pwd == "" {
		return false, ""
	}

	u := this.users[name]
	fmt.Println(u)

	if u != nil && u.PassWord == pwd {
		return true, u.Role
	}

	fmt.Println("validate failed", name, pwd)
	return false, ""
}

func LoadFromYaml(v interface{}, file string) {
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}

	err = yaml.Unmarshal(configFile, v)
	if err != nil {
		log.Fatalf("config file read error %v", err)
	}
}
