package filter

/**
 header修改
移除header、修改header内容,增加header内容
*/

//删除header中的key，可以是多个key
type HeaderRemover struct {
	Names []string
}

func (this *HeaderRemover) Name() string {
	return "header_remover"
}

func (this *HeaderRemover) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	for _, name := range this.Names {
		w.RemoveHeader(name)
	}
	return nil
}

type HeaderModify struct {
}

//新增header，可以是固定的或从context中读取的
type HeaderAdd struct {
	Key         string
	Value       string
	FromContext bool
}

func (this *HeaderAdd) Name() string {
	return "header_adder"
}

func (this *HeaderAdd) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	v := this.Value
	if this.FromContext {
		v = context[this.Value].(string)
	}

	w.SetHeader(this.Key, v)

	return nil
}

type HeaderExtract struct {
}

func (this *HeaderExtract) Name() string {
	return "header_adder"
}

func (this *HeaderExtract) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	v := "test"
	context["_test_"] = v

	return nil
}
