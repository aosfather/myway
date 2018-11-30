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
	RoundRobin LoadBalance = 0
	IPHash     LoadBalance = 1
)

//熔断
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
