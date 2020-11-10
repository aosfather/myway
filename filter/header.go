package filter

/**
 header修改
移除header、修改header内容,增加header内容

chain:
    name: xx接口处理流程
    - filter: header_remover
      parameters:
        - {source: target: config:}
*/

//删除header中的key，可以是多个key
type HeaderRemover struct {
	Parameters []SourceTarget
}

func (this *HeaderRemover) Name() string {
	return "header_remover"
}

func (this *HeaderRemover) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	for _, name := range this.Parameters {
		w.RemoveHeader(name.Source)
	}
	return nil
}

type HeaderModify struct {
}

//新增header，可以是固定的或从context中读取的
type HeaderAdd struct {
	Parameters []SourceTarget
}

func (this *HeaderAdd) Name() string {
	return "header_adder"
}

func (this *HeaderAdd) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	for _, st := range this.Parameters {
		fromContext := st.Config.(bool)
		v := st.Target
		if fromContext {
			v = context[st.Target].(string)
		}
		w.SetHeader(st.Source, v)
	}
	return nil
}
