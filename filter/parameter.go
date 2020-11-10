package filter

import (
	"fmt"
	"time"
)

/**
  参数修改
  modify[参数修改]： 参数改名、删除参数、参数移动、参数内容修改
  add[参数新增]：环境变量里的内容：例如用户名等、当前日期、当前时间、header内容
*/
//变量来源类型
type SourceStyle string

const (
	SS_CONTEXT SourceStyle = "context" //来源于上下文
	SS_DATE    SourceStyle = "date"    //来源于当前日期，其实也可以是上下文中
	SS_HEADER  SourceStyle = "header"  //来源于请求头
)

type SourceTarget struct {
	Source string
	Target string
	Config interface{}
}

//从context中读取变量，然后设置到目标参数中
//默认从context中读取
type ParameterAdder struct {
	Parameters []SourceTarget
}

func (this *ParameterAdder) Name() string {
	return "parameter_adder"
}

func (this *ParameterAdder) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	for _, st := range this.Parameters {
		style := st.Config.(SourceStyle)
		var value string
		switch style {
		case SS_HEADER:
			value = r.GetHeader(st.Source)
		case SS_CONTEXT:
			value = context[st.Source].(string)
		case SS_DATE:
			now := time.Now()
			value = now.Format(st.Source)
		}

		w.SetParamter(st.Target, value)

	}
	return nil
}

//删除参数中的key，可以是多个key
type ParamterRemover struct {
	Parameters []SourceTarget
}

func (this *ParamterRemover) Name() string {
	return "parameter_remover"
}

func (this *ParamterRemover) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	for _, name := range this.Parameters {
		w.RemoveParamter(name.Source)
	}
	return nil
}

//参数重命名
type ParamterRename struct {
	Parameters []SourceTarget
}

func (this *ParamterRename) Name() string {
	return "parameter_rename"
}

func (this *ParamterRename) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	fmt.Println(this.Parameters)
	for _, st := range this.Parameters {
		w.RemoveParamter(st.Source)
		value := r.GetParamter(st.Source)
		w.SetParamter(st.Target, value)
	}
	return nil
}
