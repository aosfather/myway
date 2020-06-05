package jobs

type TaskResult struct {
	Id     string //任务ID
	Status string //任务状态
	Msg    string //消息
	Data   map[string]string
}

//任务执行器
type TaskExecutor interface {
	Execute(taskId string, paramters map[string]string, handle func(TaskResult)) error
}

type TaskInstanceStore interface {
	UpdateTaskStatus(instanceId string, stage string, taskId string, status string)
}

//任务分配器
type TaskDispatch struct {
	store     TaskInstanceStore
	registers map[string]TaskExecutor
}

func (this *TaskDispatch) Init(store TaskInstanceStore) {
	this.store = store
	this.registers = make(map[string]TaskExecutor)
}

func (this *TaskDispatch) Register(name string, executor TaskExecutor) {
	if name != "" && executor != nil {
		this.registers[name] = executor
	}

}

func (this *TaskDispatch) Dispatch(instanceId string, stage string, t Task, context JobContext) error {
	name := t.Handle
	if te, ok := this.registers[name]; ok {
		//生成唯一的taskID，
		id := this.createTaskId(instanceId, stage, t)
		//根据input转换参数
		input := this.converToInputParameters(t, context)
		e := te.Execute(id, input, this.handleTaskResult)
		if e != nil {
			//写入错误状态
			this.store.UpdateTaskStatus(instanceId, stage, id, "error")
			return e
		}
		//写入运行状态
		this.store.UpdateTaskStatus(instanceId, stage, id, "running")
	}

	return nil
}

func (this *TaskDispatch) createTaskId(instanceId string, stage string, t Task) string {
	return ""
}

func (this *TaskDispatch) converToInputParameters(t Task, context JobContext) map[string]string {
	return nil
}

func (this *TaskDispatch) handleTaskResult(t TaskResult) {

}
