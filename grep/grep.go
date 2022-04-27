package grep

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"golang.org/x/tools/godoc/util"
)

type Config struct {
	numberWorkers int
	bufferSize    int
}

var Conf Config

const DefaultNumberWorkers = 10
const DefaultBufferSize = 100

func init() {
	Conf = Config{numberWorkers: DefaultNumberWorkers, bufferSize: DefaultBufferSize}
}

func Grep(searchPath string, searchPattern string) (*Results, error) {

	if len(searchPath) == 0 || len(searchPattern) == 0 {
		return nil, fmt.Errorf("search path and pattern must not be empty")
	}

	stat, err := os.Stat(searchPath)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return grepDir(searchPath, searchPattern)
	} else {
		result, err := GrepFile(searchPath, searchPattern)
		if err != nil {
			return nil, err
		}
		results := NewResults(Conf.bufferSize)
		results.Send(*result)
		results.Close()
		return results, nil
	}

}

func GrepFile(filepath string, searchPattern string) (*Result, error) {

	isText, err := isTextFile(filepath)
	if err != nil {
		return nil, err
	}

	result := NewResult(filepath, searchPattern)

	if !isText {
		result.IsTextFile = false
	} else {
		file, _ := os.Open(filepath)
		defer file.Close()
		regex, _ := regexp.Compile(searchPattern)
		scanner := bufio.NewScanner(file)
		lineNumber := 1
		for scanner.Scan() {
			line := scanner.Text()
			indices := regex.FindAllStringSubmatchIndex(line, -1)
			if len(indices) != 0 {
				result.AddLine(line, lineNumber, indices)
			}
			lineNumber += 1
		}
	}

	return result, nil
}

func grepDir(searchPath string, searchPattern string) (*Results, error) {

	jobs := newJobs(Conf.bufferSize)
	results := NewResults(Conf.bufferSize)

	go sendJobs(searchPath, searchPattern, jobs)
	go processJobs(jobs, results)

	return results, nil
}

func sendJobs(searchPath string, searchPattern string, jobs *Jobs) {
	sendJobsHelper(searchPath, searchPattern, jobs)
	jobs.close()
}

func sendJobsHelper(searchPath string, searchPattern string, jobs *Jobs) {

	files, err := os.ReadDir(searchPath)
	if err != nil {
		return
	}

	for _, path := range files {
		if path.IsDir() {
			dirPath := filepath.Join(searchPath, path.Name())
			sendJobsHelper(dirPath, searchPattern, jobs)
		} else {
			filePath := filepath.Join(searchPath, path.Name())
			jobs.send(filePath, searchPattern)
		}
	}
}

func processJobs(jobs *Jobs, results *Results) {

	var wg sync.WaitGroup
	for i := 0; i < Conf.numberWorkers; i++ {
		wg.Add(1)
		go deployWorker(jobs, results, &wg)
	}

	wg.Wait()
	results.Close()
}

func deployWorker(jobs *Jobs, results *Results, wg *sync.WaitGroup) {

	for job := range jobs.channel {
		result, _ := GrepFile(job.filepath, job.searchPattern)
		results.Send(*result)
	}

	wg.Done()
}

func isTextFile(filepath string) (bool, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return false, err
	}

	defer file.Close()

	data := make([]byte, 1024)
	num_bytes, err := file.Read(data)
	if err != nil {
		return false, err
	}

	return util.IsText(data[:num_bytes]), nil
}
