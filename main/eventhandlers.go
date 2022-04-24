package main

import (
	"fmt"
	"go-grep/grep"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func radioGroupHandler(selectedOption string) {
	state.UserInput.SetSelectedOption(selectedOption)
	log.Println(fmt.Sprintf("selected %v option", state.UserInput.GetSelectedOption().Value()))
}

func browseButtonHandler() {

	switch state.UserInput.GetSelectedOption() {
	case File:
		d := dialog.NewFileOpen(fileDialogButtonHandler, window)
		d.Resize(fyne.NewSize(800, 600))
		d.Show()
	case Directory:
		d := dialog.NewFolderOpen(folderDialogButtonHandler, window)
		d.Resize(fyne.NewSize(800, 600))
		d.Show()
	default:
		log.Println("Unrecognized radio option")
		window.Close()
	}
}

func formSubmitHandler() {

	infoDialog := dialog.NewInformation("Info", "Processing...", window)
	infoDialog.Show()
	err := grepAndDisplay()
	if err != nil {
		infoDialog.Hide()
		log.Println(err.Error())
		dialog.ShowError(err, window)
	} else {
		infoDialog.Hide()
	}
}

func formCancelHandler() {
	// TODO: introduction or instructions on text grid
	cwd, _ := os.Getwd()
	state.View.Clear()
	state.View.AppendText("")
	state.View.AppendText(fmt.Sprintf("Current working directory: %v", cwd))
	state.UserInput.ClearSearchPath()
	state.UserInput.ClearSearchTerm()

	log.Println("clear all")
}

func fileDialogButtonHandler(reader fyne.URIReadCloser, err error) {
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	if reader == nil {
		log.Println("file dialog cancelled")
		return
	}

	filepath := getPath(reader.URI().String())
	state.UserInput.SetSearchPath(filepath)
}

func folderDialogButtonHandler(list fyne.ListableURI, err error) {
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	if list == nil {
		log.Println("folder dialog cancelled")
		return
	}

	//children, err := list.List()
	//if err != nil {
	//	dialog.ShowError(err, window)
	//	return
	//}

	//out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
	//dialog.ShowInformation("Folder Open", out, window)

	dirpath := getPath(list.String())
	state.UserInput.SetSearchPath(dirpath)
}

func getPath(uri string) string {
	return strings.Split(uri, "file://")[1]
}

func grepAndDisplay() error {

	searchPath := state.UserInput.GetSearchPath()
	searchTerm := state.UserInput.GetSearchTerm()

	results, err := grep.Grep(searchPath, searchTerm)
	if err != nil {
		return err
	}

	state.View.Clear()
	state.View.TextGrid.ShowLineNumbers = true

	for result := range results.Channel {
		displayResult(result)
	}

	log.Println("Finished processing all results")

	return nil
}

func displayResult(result grep.Result) {
	for _, line := range result.Lines {
		state.View.AppendResult(result.Filepath, line.Line, line.LineNumber, line.Indices)
	}
}
