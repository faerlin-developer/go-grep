package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func radioGroupHandler(selectedOption string) {
	state.SetSelectedOption(selectedOption)
	log.Println(fmt.Sprintf("selected %v option", state.GetSelectedOption().Value()))
}

func browseButtonHandler() {
	switch state.GetSelectedOption() {
	case File:
		dialog.ShowFileOpen(fileDialogButtonHandler, window)
	case Directory:
		dialog.ShowFolderOpen(folderDialogButtonHandler, window)
	default:
		log.Println("Unrecognized radio option")
		window.Close()
	}
}

func formSubmitHandler() {
	log.Println("Form submitted:", state.GetSearchPath())
	log.Println("Form submitted:", state.GetSearchTerm())
	state.AppendData("new item")

	// err := errors.New("a dummy error message")
	// dialog.ShowError(err, window)
	d := dialog.NewInformation("Information", "You should know this thing...", window)
	d.Show()

	//time.Sleep(2 * time.Second)
	//d.Hide()
}

func formCancelHandler() {
	state.ClearData()
	state.ClearSearchPath()
	state.ClearSearchTerm()

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
	log.Println(reader.URI().String())

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

	children, err := list.List()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
	dialog.ShowInformation("Folder Open", out, window)
}
