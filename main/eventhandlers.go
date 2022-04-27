package main

import (
	"fmt"
	"go-grep/grep"
	"log"
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
	num_files, num_text_files, num_matched_files, num_matched_lines, err := grepAndDisplay()
	if err != nil {
		infoDialog.Hide()
		log.Println(err.Error())
		dialog.ShowError(err, window)
	} else {
		infoDialog.Hide()
		msgLine1 := fmt.Sprintf("Found %v matching lines in %v text files", num_matched_lines, num_matched_files)
		msgLine2 := fmt.Sprintf("Total number of text files found: %v", num_text_files)
		msgLine3 := fmt.Sprintf("Total number of files (including non text-files) found: %v", num_files)
		dialog.ShowInformation("Info", fmt.Sprintf("%v\n\n%v\n%v", msgLine1, msgLine2, msgLine3), window)
	}

	if num_matched_lines == 0 {
		state.View.ShowDefaultView()
	}

	log.Println(state.UserInput.GetNumberWorkers())
	log.Println(state.UserInput.GetBufferSize())

	log.Println(num_text_files)
	log.Println(num_matched_lines)

}

func formCancelHandler() {
	state.View.Clear()
	state.View.ShowDefaultView()
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

	dirpath := getPath(list.String())
	state.UserInput.SetSearchPath(dirpath)
}

func getPath(uri string) string {
	return strings.Split(uri, "file://")[1]
}

func grepAndDisplay() (int, int, int, int, error) {

	searchPath := state.UserInput.GetSearchPath()
	searchTerm := state.UserInput.GetSearchTerm()

	results, err := grep.Grep(searchPath, searchTerm)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	state.View.Clear()
	state.View.TextGrid.ShowLineNumbers = true

	num_files := 0
	num_text_files := 0
	num_matched_files := 0
	num_matched_lines := 0

	for result := range results.Channel {
		num_files += 1
		if result.IsTextFile {
			num_text_files += 1
		}
		if !result.IsEmpty() {
			displayResult(result)
			num_matched_lines += len(result.Lines)
			num_matched_files += 1
		}
	}
	log.Println("Finished processing results")

	return num_files, num_text_files, num_matched_files, num_matched_lines, nil
}

func displayResult(result grep.Result) {
	for _, line := range result.Lines {
		state.View.AppendResult(result.Filepath, line.Line, line.LineNumber, line.Indices)
	}
}
