package jobs

type TaskInstanceStore interface {
	UpdateTaskStatus(id string, status string, msg string)
}

//job状态存储
type JobInstanceStore interface {
	LoadJobInstance(instanceId string) *JobInstance
	UpdateCurrentStage(instanceId string, code string)
	CreateInstance(instanceId string, job string)
	SaveContext(instanceId string, context JobContext)
}

type MetaManager interface {
	AddJob(job Job)
	GetJob(j string) Job
	GetTask(t string) Task
	AddTask(t Task)
}
