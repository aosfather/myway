package jobs

//job上下文接口
type JobContext interface {
	GetContext(key string) string
	SetContext(key string, v string)
	Visitor(func(key, v string))
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
