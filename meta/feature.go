package meta

/**
  特性
*/

//查找规则
type MatchRule int32

const (
	MatchDefault MatchRule = 0
	MatchAll     MatchRule = 1
	MatchAny     MatchRule = 2
)

// LoadBalance the load balance enum
type LoadBalance int32

const (
	LBRoundRobin LoadBalance = 0
	LBIPHash     LoadBalance = 1
	LBTag        LoadBalance = 2
)

/**
熔断器
circuitBreaker.errorThresholdPercentage
失败率达到多少百分比后熔断
默认值：50
主要根据依赖重要性进行调整

circuitBreaker.forceClosed
是否强制关闭熔断
如果是强依赖，应该设置为true

circuitBreaker.requestVolumeThreshold
熔断触发的最小个数/10s
默认值：20

circuitBreaker.sleepWindowInMilliseconds
熔断多少秒后去尝试请求
默认值：5000

hystrix.command.default.metrics.healthSnapshot.intervalInMilliseconds
记录health 快照（用来统计成功和错误绿）的间隔，默认500ms
*/
type CircuitBreaker struct {
	CloseTimeout       int64 `protobuf:"varint,1,opt,name=closeTimeout" json:"closeTimeout"`
	HalfTrafficRate    int32 `protobuf:"varint,2,opt,name=halfTrafficRate" json:"halfTrafficRate"`
	RateCheckPeriod    int64 `protobuf:"varint,3,opt,name=rateCheckPeriod" json:"rateCheckPeriod"`
	FailureRateToClose int32 `protobuf:"varint,4,opt,name=failureRateToClose" json:"failureRateToClose"`
	SucceedRateToOpen  int32 `protobuf:"varint,5,opt,name=succeedRateToOpen" json:"succeedRateToOpen"`
}

//熔断状态
type CircuitStatus int32

const (
	CS_Open  CircuitStatus = 0
	CS_Half  CircuitStatus = 1
	CS_Close CircuitStatus = 2
)
