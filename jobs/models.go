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

//任务
type Task struct {
	Label     string            //任务名称
	Handle    string            //任务处理器code
	InputMap  map[string]string //输入参数mapping key 参数key value 输入参数key
	OutputMap map[string]string //输出参数mapping key out参数可以，value 放入context中的参数key
}

//job上下文接口
type JobContext interface {
	GetContext(key string) string
	SetContext(key string, v string)
	Visitor(func(key, v string))
}

type JobStore interface {
	LoadJobInstance(instanceId string) *JobInstance
}

//job状态存储
type JobInstanceStore interface {
	UpdateCurrentStage(instanceId string, code string)
	SaveContext(instanceId string, context JobContext)
}

//步骤执行器
type StageExecutor func(stage Stage, context JobContext) (string, error)

//单个job实例
type JobInstance struct {
	job          *Job
	context      JobContext
	currentStage string           //当前步骤
	store        JobInstanceStore //持久化
}

func (this *JobInstance) Next() bool {
	return false
}

func (this *JobInstance) executeStage(s Stage) {

}
