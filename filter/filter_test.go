package filter

import (
	"github.com/aosfather/myway/runtime"
	"testing"
)

type testReader struct {
}

func (this *testReader) GetHeader(key string) string {
	return "header1"
}

func (this *testReader) GetParamter(key string) string {
	return "parameter1"
}

func (this *testReader) GetBody() []byte {
	return []byte("parameter1=1")
}

func TestFilterChain_DoFilter(t *testing.T) {
	context := make(runtime.RunTimeContext)
	t.Log(context)
}

func TestFilterChain_factory(t *testing.T) {
	context := make(runtime.RunTimeContext)
	t.Log(context)
	fdef := FilterDef{Filter: "header_remover", Parameters: []SourceTarget{SourceTarget{Source: "v"}}}
	cdef := FilterChainDef{Name: "test", Filters: []FilterDef{fdef}}
	c := Factory(cdef)
	//c.DoFilter(nil,nil,context)
	t.Log(c[0])
}
