package main

import (
	"go-grep/grep"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var state State
var window fyne.Window

func init() {
	state = NewState()
	window = NewWindow()
}

func main() {

	result := grep.Grep("test_dir/file1.txt", "fox")
	for _, line := range result.Line() {
		log.Println(line)
	}

	top := createTopComponent()
	middle := createMiddleComponent()
	bottom := createBottomComponent()
	content := container.New(layout.NewBorderLayout(top, bottom, nil, nil),
		top, bottom, middle)

	window.SetContent(content)
	window.Resize(fyne.NewSize(900, 700))
	window.ShowAndRun()

}
