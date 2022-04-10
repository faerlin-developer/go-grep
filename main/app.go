package main

import (
	"fmt"

	"go-grep/grep"
)

func main() {

	result := grep.Grep("test_dir/file1.txt", "fox")
	for _, line := range result.Line() {
		fmt.Println(line)
	}

	Initialize()

}
