package runtime

import (
	"fmt"
	"testing"
	"time"
)

func TestRuntimeContext_QPSCount(t *testing.T) {
	context := runtimeContext{}
	context.Init()
	context.QPS.Max = 50
	for i := 0; i < 200; i++ {
		fmt.Println("-", i)
		context.QPS.Incr()
		time.Sleep(10 * time.Millisecond)
		//fmt.Println(*context.QPS)
	}

	time.Sleep(1 * time.Second)
}
