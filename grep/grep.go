package grep

import (
	"bufio"
	"os"
	"strings"
)

func Grep(filepath string, searchTerm string) *Result {

	file, err := os.Open(filepath)
	if err != nil {
		return nil
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

	if result.isEmpty() {
		return nil
	} else {
		return result
	}

}
