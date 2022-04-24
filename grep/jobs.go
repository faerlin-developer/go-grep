package grep

type Job struct {
	filepath      string
	searchPattern string
}

type Jobs struct {
	channel chan Job
}

func (jobs *Jobs) send(filepath string, searchPattern string) {
	jobs.channel <- Job{filepath, searchPattern}
}

func (jobs *Jobs) close() {
	close(jobs.channel)
}

func newJobs(size int) *Jobs {
	jobs := make(chan Job, size)
	return &Jobs{jobs}
}
