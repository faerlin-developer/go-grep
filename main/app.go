package main

import (
	"fmt"

	"go-grep/grep"
	"go-grep/ui"
)

func main() {

	result := grep.Grep("test_dir/file1.txt", "fox")
	for _, line := range result.Line() {
		fmt.Println(line)
	}

	ui.Initialize()

}
