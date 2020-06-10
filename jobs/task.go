package jobs

//任务实例
type TaskInstance struct {
	Status TaskStatus
	Task   *Task  //任务元信息
	Id     string //实例ID
	Input  map[string]string
}

func (this *TaskInstance) buildInput(context JobContext) {
	if this.Task == nil {
		return
	}

	if this.Task.InputMap != nil && len(this.Task.InputMap) > 0 {
		this.Input = make(map[string]string)
		for key, value := range this.Task.InputMap {
			if value.Refer {
				v := context.GetContext(value.Value)
				this.Input[key] = v
			} else {
				this.Input[key] = value.Value
			}

		}
	}
}

//任务状态
type TaskStatus struct {
	Label      string
	Terminal   bool
	Successful bool
	Retriable  bool
}

var (
	_IN_PROGRESS                = TaskStatus{"inProgress", false, true, true}
	_CANCELED                   = TaskStatus{"canceled", true, false, false}
	_FAILED                     = TaskStatus{"failed", true, false, true}
	_FAILED_WITH_TERMINAL_ERROR = TaskStatus{"failedWithTerminalError", true, false, false}
	_COMPLETED                  = TaskStatus{"completed", true, true, true}
	_COMPLETED_WITH_ERRORS      = TaskStatus{"completedWithErrors", true, true, true}
	_SCHEDULED                  = TaskStatus{"scheduled", false, true, true}
	_TIMED_OUT                  = TaskStatus{"timedOut", true, false, true}
	_SKIPPED                    = TaskStatus{"skipped", true, true, false}
)

//任务执行状态
type TaskExecuteStatus byte

const (
	TS_IN_PROGRESS                TaskExecuteStatus = 1 //执行中
	TS_FAILED                     TaskExecuteStatus = 2 //失败
	TS_FAILED_WITH_TERMINAL_ERROR TaskExecuteStatus = 3
	TS_COMPLETED                  TaskExecuteStatus = 5 //完成
)

type TaskExecLog struct {
	TaskId      string //任务ID
	CreatedTime int64  //创建时间
	Content     string //日志内容
}

type TaskResult struct {
	Id            string            //任务实例ID
	TaskInstance  *TaskInstance     //任务定义的code
	FlowId        string            //流程实例ID
	WorkerId      string            //执行者
	Logs          []TaskExecLog     //执行日志
	SubWorkFlowId string            //子流程ID
	Status        TaskExecuteStatus //任务状态
	Msg           string            //消息
	Data          map[string]string
}

func (this *TaskResult) PutToContext(context JobContext) {
	if this.TaskInstance != nil && context != nil {
		outputmap := this.TaskInstance.Task.OutputMap
		if outputmap != nil && len(outputmap) > 0 {
			for key, value := range outputmap {
				v := this.Data[key]
				context.SetContext(value, v)
			}
		}

	}

}

//任务变更监听
type TaskChangeListener interface {
	TaskChangeEventNotify(result *TaskResult)
}
