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
	tf := HeaderExtract{}
	tf.DoFilter(nil, nil, context)
	t.Log(context)
}
