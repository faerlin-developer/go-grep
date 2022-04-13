package grep

type Job struct {
	searchTerm string
	filepath   string
}

type JobChannel struct {
	jobs chan Job
}

func (jc *JobChannel) Send(filepath string, searchTerm string) {
	jc.jobs <- Job{filepath, searchTerm}
}

func (jc *JobChannel) Receive() Job {
	return <-jc.jobs
}

func NewJobChannel(size int) *JobChannel {
	jobs := make(chan Job, size)
	return &JobChannel{jobs}
}
