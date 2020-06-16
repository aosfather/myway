package jobs

import (
	"fmt"
)

//任务执行器
type TaskExecutor interface {
	//执行
	Execute(task TaskInstance, paramters map[string]string, handle func(*TaskResult)) error
}

//任务分配器
type TaskDispatch struct {
	TaskMan   MetaManager
	store     TaskInstanceStore
	registers map[string]TaskExecutor
	listeners []TaskChangeListener
}

func (this *TaskDispatch) Init(store TaskInstanceStore) {
	this.store = store
	this.registers = make(map[string]TaskExecutor)

}

func (this *TaskDispatch) Addlistener(t ...TaskChangeListener) {
	if len(t) > 0 {
		this.listeners = append(this.listeners, t...)
	}
}

func (this *TaskDispatch) Register(name string, executor TaskExecutor) {
	if name != "" && executor != nil {
		this.registers[name] = executor
	}

}

func (this *TaskDispatch) Dispatch(t TaskInstance, context JobContext) error {
	name := t.Task.Code
	if te, ok := this.registers[name]; ok {
		t.buildInput(context)
		e := te.Execute(t, t.Input, this.handleTaskResult)
		if e != nil {
			//写入错误状态
			this.store.UpdateTaskStatus(t.Id, "error", e.Error())
			return e
		}
		//写入运行状态
		this.store.UpdateTaskStatus(t.Id, "running", "ok")
	}

	return nil
}

func (this *TaskDispatch) CreateTaskId(instanceId string, stage string, t Task) string {
	return fmt.Sprintf("%s_%s_%s", instanceId, stage, t.Code)
}

//转换输出参数
func (this *TaskDispatch) converFromOutputParameters(t Task, out map[string]string) map[string]string {
	if t.OutputMap != nil && len(t.OutputMap) > 0 {
		context := make(map[string]string)
		for key, value := range t.OutputMap {
			v := out[key]
			context[value] = v
		}

		return context
	}

	return nil

}

func (this *TaskDispatch) handleTaskResult(t *TaskResult) {
	//如果失败，进入重试

	//如果确认失败，进入补偿
	if t.Status == TS_COMPLETED {
		//更新记录
		this.store.UpdateTaskStatus(t.Id, "", "")
		//通知相关listener
		if this.listeners != nil {
			for _, l := range this.listeners {
				if l != nil {
					l.TaskChangeEventNotify(t)
				}

			}
		}
	}

}
