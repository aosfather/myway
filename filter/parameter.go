package filter

import (
	"time"
)

/**
  参数修改
  modify[参数修改]： 参数改名、删除参数、参数移动、参数内容修改
  add[参数新增]：环境变量里的内容：例如用户名等、当前日期、当前时间、header内容
*/
//变量来源类型
type SourceStyle byte

const (
	SS_CONTEXT SourceStyle = 0 //来源于
	SS_DATE    SourceStyle = 1
	SS_HEADER  SourceStyle = 4
)

//从context中读取变量，然后设置到目标参数中
type ParameterAdder struct {
	Source string
	Target string
	Style  SourceStyle //是否从header中读取，默认从context中读取
}

func (this *ParameterAdder) Name() string {
	return "parameter_adder"
}

func (this *ParameterAdder) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	var value string
	switch this.Style {
	case SS_HEADER:
		value = r.GetHeader(this.Source)
	case SS_CONTEXT:
		value = context[this.Source].(string)
	case SS_DATE:
		now := time.Now()
		value = now.Format(this.Source)
	}

	w.SetParamter(this.Target, value)
	return nil
}

//删除参数中的key，可以是多个key
type ParamterRemover struct {
	Names []string
}

func (this *ParamterRemover) Name() string {
	return "parameter_remover"
}

func (this *ParamterRemover) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	for _, name := range this.Names {
		w.RemoveParamter(name)
	}
	return nil
}

//参数重命名
type ParamterRename struct {
	Source string
	Target string
}

func (this *ParamterRename) Name() string {
	return "parameter_rename"
}

func (this *ParamterRename) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	w.RemoveParamter(this.Source)
	value := r.GetParamter(this.Source)
	w.SetParamter(this.Target, value)

	return nil
}
