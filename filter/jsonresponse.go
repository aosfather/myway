package filter

import (
	"encoding/json"
	"fmt"
	"github.com/aosfather/bingo_utils/codes"
	"strings"
)

/**
  json 返回参数的处理
*/
type BodyToJson struct {
}

func (this *BodyToJson) Name() string {
	return "body_json"
}

func (this *BodyToJson) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	body := r.GetBody()
	response := make(map[string]interface{})
	json.Unmarshal(body, &response)
	context["_response_"] = response
	return nil
}

//将对象转换成json写入body中
type JsonToBody struct {
}

func (this *JsonToBody) Name() string {
	return "json_body"
}

func (this *JsonToBody) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	body := context["_response_"]
	if body != nil {
		jsondata, err := json.Marshal(body)
		w.SetBody(jsondata)
		return err
	}
	return nil
}

//属性处理函数
type PropertyProcess func(obj map[string]interface{}, source, target string)

//通用属性处理filter
type PropertyFilter struct {
	Parameters []SourceTarget
	name       string
	process    PropertyProcess
}

func (this *PropertyFilter) Name() string {
	return this.name
}

func (this *PropertyFilter) DoFilter(r EntityReader, w EntityWriter, context map[string]interface{}) error {
	body := context["_response_"].(map[string]interface{})
	if body != nil {
		for _, st := range this.Parameters {
			objs, p := getObject(body, st.Source)
			for _, obj := range objs {
				this.process(obj, p, st.Target)
			}
		}
	}
	return nil
}

//修改属性名称
func propertyRenname(obj map[string]interface{}, source, target string) {
	if obj == nil {
		return
	}
	v := obj[source]
	delete(obj, source)
	obj[target] = v
}

//脱敏处理
func propertyDesensitation(obj map[string]interface{}, source, target string) {
	if obj == nil {
		return
	}
	v := obj[source]
	if v != nil {
		desensitation := desensitatinMap[target]
		if desensitation != nil {
			obj[source] = desensitation.Convert(v.(string))
		}
	}

}

var desensitatinMap = make(map[string]*codes.Sensitive)

//注册使用的脱敏处理类型
func RegisterDesensitive(name string, desensitation *codes.Sensitive) {
	if name != "" && desensitation != nil {
		desensitatinMap[name] = desensitation
	}
}

//根据path获取父级别对象。仅仅支持末级别为数组的形式，对于path中间为数组的不支持。
func getObject(body map[string]interface{}, source string) ([]map[string]interface{}, string) {
	paths := strings.Split(source, "/")
	if len(paths) <= 1 {
		return []map[string]interface{}{body}, source
	}
	size := len(paths) - 1
	theParentPaths := paths[:size]
	p := paths[size]
	obj := body
	//获取操作的对象，逐层的处理
	for level, key := range theParentPaths {
		fmt.Println(level)
		v := obj[key]
		if v != nil {
			if objs, ok := v.([]interface{}); ok {
				var nobjs []map[string]interface{}
				for _, o := range objs {
					nobjs = append(nobjs, o.(map[string]interface{}))
				}
				return nobjs, p
			}
			obj = v.(map[string]interface{})
		} else {
			obj = nil
		}

	}

	return []map[string]interface{}{obj}, p
}
