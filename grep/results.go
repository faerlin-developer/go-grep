package grep

type Line struct {
	Line       string
	LineNumber int
	Indices    [][]int
}

type Result struct {
	Filepath      string
	SearchPattern string
	Lines         []Line
}

type Results struct {
	Channel chan Result
}

func (r *Result) AddLine(line string, lineNumber int, indices [][]int) {
	r.Lines = append(r.Lines, Line{line, lineNumber, indices})
}

func (r *Result) IsEmpty() bool {
	return len(r.Lines) == 0
}

func (results *Results) Send(result Result) {
	results.Channel <- result
}

func (results *Results) Close() {
	close(results.Channel)
}

func NewResult(filepath string, searchPattern string) *Result {
	return &Result{filepath, searchPattern, make([]Line, 0)}
}

func NewResults(size int) *Results {
	results := make(chan Result, size)
	return &Results{results}
}
