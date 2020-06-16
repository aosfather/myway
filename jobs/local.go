package jobs

/**
  本地实现
*/
var (
	_NilJob  = Job{}
	_NilTask = Task{}
)

type record struct {
	Job   string
	Stage string
}
type LocalMetaManager struct {
	jobs         map[string]Job
	tasks        map[string]Task
	jobInstances map[string]record
	jobContexts  map[string]JobContext
}

func (this *LocalMetaManager) Init() {
	this.jobs = make(map[string]Job)
	this.tasks = make(map[string]Task)
	this.jobContexts = make(map[string]JobContext)
	this.jobInstances = make(map[string]record)
}

func (this *LocalMetaManager) GetTask(t string) Task {
	if task, ok := this.tasks[t]; ok {
		return task
	}

	return _NilTask
}

func (this *LocalMetaManager) AddTask(t Task) {
	if _, ok := this.tasks[t.Code]; !ok {
		this.tasks[t.Code] = t
	}
}

func (this *LocalMetaManager) AddJob(job Job) {
	if _, ok := this.jobs[job.Code]; !ok {
		this.jobs[job.Code] = job
	}
}

func (this *LocalMetaManager) GetJob(j string) Job {
	job, ok := this.jobs[j]
	if ok {
		return job
	}
	return _NilJob
}

func (this *LocalMetaManager) LoadJobInstance(instanceId string) *JobInstance {
	code, ok := this.jobInstances[instanceId]
	if ok {
		job := this.GetJob("")
		if job.Code == "" {
			return nil
		}
		instance := JobInstance{&job, this.jobContexts[instanceId], code, this}
		return &instance
	}

	return nil
}

func (this *LocalMetaManager) CreateInstance(instanceId string, job string) {
	if instanceId == "" || job == "" {
		return
	}
	_, ok := this.jobInstances[instanceId]
	if ok {
		return
	} else {
		this.jobInstances[instanceId] = record{Job: job}
	}
}

func (this *LocalMetaManager) UpdateCurrentStage(instanceId string, code string) {
	if instanceId == "" || code == "" {
		return
	}
	if v, ok := this.jobInstances[instanceId]; ok {
		v.Stage = code
	}
}

func (this *LocalMetaManager) SaveContext(instanceId string, context JobContext) {
	if instanceId == "" || context == nil {
		return
	}
	this.jobContexts[instanceId] = context
}
