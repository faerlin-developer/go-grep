package grep

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type Config struct {
	numberWorkers int
	bufferSize    int
}

var Conf Config

func init() {
	Conf = Config{numberWorkers: 10, bufferSize: 100}
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

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	regex, _ := regexp.Compile(searchPattern)
	result := NewResult(filepath, searchPattern)
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
		if !result.IsEmpty() {
			results.Send(*result)
		}
	}

	wg.Done()
}
