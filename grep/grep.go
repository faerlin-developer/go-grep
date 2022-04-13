package grep

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var Conf Config

func init() {
	Conf = Config{numberWorkers: 10, bufferSize: 100}
}

func GrepFile(filepath string, searchTerm string) (*Result, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	result := NewResult(filepath, searchTerm)
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchTerm) {
			result.AddLine(line, lineNumber)
		}
		lineNumber += 1
	}

	return result, nil
}

func Grep(searchPath string, searchTerm string) (*ResultChannel, error) {

	jobChannel := NewJobChannel(Conf.GetBufferSize())
	resultChannel := NewResultChannel(Conf.GetBufferSize())

	go sendJobs(searchPath, searchTerm, jobChannel)
	go processJobs(jobChannel, resultChannel)

	return resultChannel, nil

}

func sendJobs(searchPath string, searchTerm string, jobChannel *JobChannel) {

	files, err := os.ReadDir(searchPath)
	if err != nil {
		return
	}

	for _, path := range files {
		if path.IsDir() {
			dirPath := filepath.Join(searchPath, path.Name())
			sendJobs(dirPath, searchTerm, jobChannel)
		} else {
			filePath := filepath.Join(searchPath, path.Name())
			jobChannel.Send(filePath, searchTerm)
		}
	}
}

func processJobs(jobChannel *JobChannel, resultChannel *ResultChannel) {
	// ...
}
