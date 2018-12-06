package runtime

import (
	"sync/atomic"
	"time"
)

/**
  运行时记录
*/

type CiruitStatus byte

const (
	CS_OPEN  CiruitStatus = 2 //熔断打开
	CS_CLOSE CiruitStatus = 0 //熔断关闭
	CS_HALF  CiruitStatus = 1 //熔断半打开状态

)

//运行时上下文
type runtimeContext struct {
	ID     string       //
	Owner  interface{}  //所属对象
	QPS    QPSCount     //访问量
	Status CiruitStatus //熔断状态
	Lb     LoadBalance  //负载均衡器

}

func (this *runtimeContext) Init() {
	this.QPS.Init()
}

//QPS记录器
type QPSCount struct {
	Max   int64
	qps   *int64
	times *int64
}

func (this *QPSCount) Init() {
	this.qps = new(int64)
	this.times = new(int64)
}

//计数
func (this QPSCount) Incr() bool {
	//先看是否在同一秒里,如果不是，则需要重新计数
	t := time.Now().Unix()
	if t != atomic.LoadInt64(this.times) {
		atomic.StoreInt64(this.times, t)
		atomic.StoreInt64(this.qps, 0)
	} else {
		//限制不允许超过max
		if atomic.LoadInt64(this.qps) >= this.Max {
			return false
		}
	}

	atomic.AddInt64(this.qps, 1)
	return true
}

func (this *QPSCount) Count() int64 {
	return atomic.LoadInt64(this.qps)
}

//
