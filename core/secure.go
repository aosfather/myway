package core

import (
	"fmt"
	"strings"
)

//安全模型
const MAX_LEVEL = 10

//元模型管理接口
type RoleMetaManager interface {
	FindRole(name string) *Role
}

type DefaultRoleMetaManager struct {
	roles map[string]*Role
}

func (this *DefaultRoleMetaManager) AddRole(r *Role) {
	if this.roles == nil {
		this.roles = make(map[string]*Role)
	}

	if r != nil {
		name := r.Code
		if this.roles[name] == nil {
			this.roles[name] = r
		}
	}
}

func (this *DefaultRoleMetaManager) FindRole(name string) *Role {
	if len(this.roles) > 0 {
		return this.roles[name]
	}
	return nil
}

//角色管理
type RoleManager struct {
	meta RoleMetaManager
}

func (this *RoleManager) SetMetaManager(m RoleMetaManager) {
	this.meta = m
}

func (this *RoleManager) Validate(rolename string, key string, obj map[string]string) bool {
	if key == "" || len(obj) == 0 {
		return false
	}

	if this.meta == nil {
		return false
	}

	r := this.meta.FindRole(rolename)
	if r != nil {
		return r.Validate(key, obj, 0)
	}

	return false
}

//角色
type Role struct {
	SuperRole *Role  //继承的角色
	Code      string //角色名称
	setMap    map[string]*AuthobjectSet
}

//增加权限集
func (this *Role) AddAuthObjectSet(set *AuthobjectSet) {
	if this.setMap == nil {
		this.setMap = make(map[string]*AuthobjectSet)
	}

	this.setMap[set.Class.Code] = set
}

func (this *Role) ValidateSelf(key string, obj map[string]string) bool {
	set := this.setMap[key]
	if set == nil {
		return false
	}

	return set.Validate(key, obj)
}

func (this *Role) Validate(key string, obj map[string]string, level int32) bool {
	if level > MAX_LEVEL {
		//超过最大链接长度,防止循环继承
		return false
	}
	//先校验自身
	if this.ValidateSelf(key, obj) {
		return true
	} else {
		//如果有继承，使用继承的角色权限进行校验
		if this.SuperRole == nil {
			return false
		}

		return this.SuperRole.Validate(key, obj, level+1)
	}

}

//权限对象集合
type AuthobjectSet struct {
	Class         *AuthClass
	Objects       []*Authobject
	ReverseObject []*Authobject
}

func (this *AuthobjectSet) BuildFieldValue(name string, t ValueType, v string) *AuthFieldValue {
	f := this.Class.Fields[name]
	if f != nil {
		return &AuthFieldValue{f, t, v, nil}
	}

	return nil
}

func (this *AuthobjectSet) AddAuthObject(reverse bool, values []*AuthFieldValue) {
	obj := Authobject{reverse, values}
	if reverse {
		this.ReverseObject = append(this.ReverseObject, &obj)
	} else {
		this.Objects = append(this.Objects, &obj)
	}

}

func (this *AuthobjectSet) Validate(key string, obj map[string]string) bool {
	if key != this.Class.Code {
		return false
	}

	//校验正向，只要一个满足则完成校验，表示通过
	for _, authObjs := range this.Objects {
		if authObjs.Validate(obj) {
			return true
		}
	}

	if len(this.ReverseObject) == 0 {
		return false
	}

	//先校验黑名单，反向，只要全部不满足则完成校验，一个满足则校验不通过
	for _, authObjs := range this.ReverseObject {
		if !authObjs.Validate(obj) {
			return false
		}
	}

	return true
}

//权限对象
type Authobject struct {
	Reverse bool //反向，取反。
	Values  []*AuthFieldValue
}

func (this *Authobject) Validate(obj map[string]string) bool {
	v := this.in(obj)
	if this.Reverse {
		return !v
	}
	return v
}

func (this *Authobject) in(obj map[string]string) bool {
	var name, value string
	for _, v := range this.Values {
		name = v.Field.Code
		value = obj[name]
		if !v.Validate(value) {
			return false
		}
	}
	return true
}

//权限值
type AuthFieldValue struct {
	Field  *AuthField //字段名称
	Type   ValueType  //取值类型：单值、区间、枚举
	Value  string     //取值范围
	Values []string
}

func (this *AuthFieldValue) Validate(v string) bool {
	switch this.Type {
	case VT_SINGLE:
		return this.Value == v
		//区间
	case VT_RANGE:
		//集合
	case VT_SET:
		if len(this.Values) == 0 {
			this.Values = strings.Split(this.Value, ",")
		}
		for _, value := range this.Values {
			if value == v {
				return true
			}
		}
	}

	return false
}

//权限类，模板
type AuthClass struct {
	Code   string                //权限类名
	Fields map[string]*AuthField //权限字段
}

func (this *AuthClass) AddFieldByParame(name, desc string, t FieldType) {
	if name == "" {
		return
	}

	f := AuthField{name, desc, t}
	this.AddField(&f)
}

func CreateAuthClass(c string) *AuthClass {
	class := AuthClass{}
	class.Code = c
	return &class
}

func (this *AuthClass) AddField(f *AuthField) {
	if f == nil || f.Code == "" {
		return
	}
	if this.Fields == nil {
		this.Fields = make(map[string]*AuthField)
	}

	this.Fields[f.Code] = f

}

//权限字段
type AuthField struct {
	Code string    //字段ID
	Desc string    //描述
	Type FieldType //类型
}

func (this *AuthField) ToString(o interface{}) string {

	return fmt.Sprintf("%s", o)
}

//值类型
type ValueType byte

const (
	VT_SINGLE ValueType = 1 //单值
	VT_RANGE  ValueType = 2 //区间
	VT_SET    ValueType = 3 //集合（枚举)
)

//字段类型
type FieldType byte

const (
	FT_STRING FieldType = 1 //字符
	FT_INT    FieldType = 2 //整数
	FT_DATE   FieldType = 3 //日期
)
