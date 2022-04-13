package grep

type Result struct {
	filepath   string
	searchTerm string
	line       []string
	lineNumber []int
}

type ResultChannel struct {
	results chan Result
}

func (r *Result) AddLine(line string, lineNumber int) {
	r.line = append(r.line, line)
	r.lineNumber = append(r.lineNumber, lineNumber)
}

func (r *Result) isEmpty() bool {
	return len(r.line) == 0
}

func (r *Result) Line() []string {
	return r.line
}

func (rc *ResultChannel) Send(result Result) {
	rc.results <- result
}

func (rc *ResultChannel) Receive() Result {
	return <-rc.results
}

func NewResult(filepath string, searchTerm string) *Result {
	return &Result{filepath, searchTerm, make([]string, 0), make([]int, 0)}
}

func NewResultChannel(size int) *ResultChannel {
	results := make(chan Result, size)
	return &ResultChannel{results}
}
