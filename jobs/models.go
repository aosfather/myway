package jobs

/**
  基本模型
 1、Job  工作，job由 stage构成
 2、stage 工作中的步骤和阶段。类型有：goto (选择判断后续的stage分支)、partitions 分片，用于执行并发任务、normal 标准常规步骤只有单个task执行
 3、task  在某个阶段执行的任务
执行期间
 JobInstance 实例
 stage间可以通过context传递消息和数据
 taskdispatch 任务调度分配
 TaskExecutor 用于执行处理单个任务
*/

//工作
type Job struct {
	Code   string  //工作编码
	Label  string  //工作名称
	stages []Stage //阶段列表

}

//工作阶段
type Stage struct {
	Code          string //阶段标识
	Label         string //阶段名称
	Type          string //类型
	PartitionTask Task   //分片任务
	WorkTask      Task   //任务
}

type Value struct {
	Refer bool   //是否引用
	Value string //取值
}

//任务类型
type TaskType byte

const (
	//基本任务
	TT_SIMPLE TaskType = 10
	//动态任务
	TT_DYNAMIC TaskType = 11
	//Fork
	TT_FORK_JOIN TaskType = 12
	//动态fork
	TT_FORK_JOIN_DYNAMIC TaskType = 13
	//选择分支
	TT_DECISION TaskType = 20
	//join
	TT_JOIN TaskType = 14
	//
	TT_EXCLUSIVE_JOIN TaskType = 15
	//循环
	TT_DO_WHILE TaskType = 21
	//子流程
	TT_SUB_WORKFLOW TaskType = 22
	//事件
	TT_EVENT TaskType = 30
	//等待
	TT_WAIT TaskType = 23
	//http访问
	TT_HTTP TaskType = 41
	//内嵌表达式
	TT_LAMBDA TaskType = 42
	//结束
	TT_TERMINATE TaskType = 24
	//广播到kafka
	TT_KAFKA_PUBLISH TaskType = 43
	//用户自定义
	TT_USER_DEFINED TaskType = 50
)

//是否系统任务
func IsSystemTask(t TaskType) bool {
	return t < 50
}

//任务
type Task struct {
	Code      string            //任务唯一编码
	Type      string            //任务类型，system、simple
	Label     string            //任务名称
	InputMap  map[string]Value  //输入参数mapping key 参数key value 输入参数key
	OutputMap map[string]string //输出参数mapping key out参数key，value 放入context中的参数key
	Retry     RetryConfig       //重试设置
	Timeout   TimeoutConfig     //超时设置
}

//重试逻辑
type RetryLogic byte

const (
	//重试
	TP_RETRY TimeoutPolicy = 1
	//流程超时
	TP_TIME_OUT_WF TimeoutPolicy = 2
	//警告
	TP_ALERT_ONLY TimeoutPolicy = 3
)

//超时策略
type TimeoutPolicy byte

const (
	//固定delay时间后重试
	RL_FIXED RetryLogic = 1
	//指数级延时后重试 :  retryDelaySeconds * attempNo 之后重新调度任务
	RL_EXPONENTIAL_BACKOFF RetryLogic = 2
)

//超时设置
type TimeoutConfig struct {
	TimeoutSeconds         int           //超时时间（秒）
	Policy                 TimeoutPolicy //超时策略
	ResponseTimeoutSeconds int           //如果大于0，则如果在此时间后未更新状态，则重新调度任务。
	PollTimeoutSeconds     int           //拉取超时设置
}

//重试设置
type RetryConfig struct {
	Count       int        //重试次数
	Logic       RetryLogic //策略
	DelaySecond int        //重试延时
}
