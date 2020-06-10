package jobs

type TaskInstanceStore interface {
	UpdateTaskStatus(id string, status string, msg string)
}

type JobStore interface {
	LoadJobInstance(instanceId string) *JobInstance
}

//job状态存储
type JobInstanceStore interface {
	UpdateCurrentStage(instanceId string, code string)
	SaveContext(instanceId string, context JobContext)
}

type TaskManager interface {
	GetTask(t string) Task
	AddTask(t Task)
}
