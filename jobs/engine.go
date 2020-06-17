package jobs

type Engine struct {
	meta     MetaManager
	jobstore JobInstanceStore
}

func (this *Engine) Start(jobId string) *JobInstance {
	if this.meta != nil {
		job := this.meta.GetJob(jobId)
		if job.Code != "" {
			this.jobstore.CreateInstance("", job.Code)
			return &JobInstance{job: &job, context: nil, currentStage: job.Stages[0].Code, store: this.jobstore}
		}
	}
	return nil
}
