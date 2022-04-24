package main

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createMiddleComponent() *container.Scroll {

	cwd, _ := os.Getwd()
	state.View.AppendText("")
	state.View.AppendText("")
	state.View.AppendText(fmt.Sprintf("Current working directory: %v", cwd))

	return state.View.Scroll
}

func createBottomComponent() *fyne.Container {

	radio := widget.NewRadioGroup([]string{File.Value(), Directory.Value()}, radioGroupHandler)
	radio.Horizontal = true
	radio.SetSelected(File.Value())

	searchPatternEntry := widget.NewEntryWithData(*state.UserInput.SearchPattern)
	searchPathEntry := widget.NewEntryWithData(*state.UserInput.SearchPath)
	searchPathEntry.PlaceHolder = "Enter absolute path or path relative to current working directory"

	browseButton := widget.NewButton("Browse", browseButtonHandler)
	browseButtonContainer := container.New(layout.NewHBoxLayout(), browseButton)
	searchPathContainer := container.New(layout.NewBorderLayout(nil, nil, nil, browseButtonContainer),
		browseButtonContainer, searchPathEntry)

	radioFormItem := widget.NewFormItem("", radio)
	searchPathFormItem := widget.NewFormItem("Search Path", searchPathContainer)
	searchPatternFormItem := widget.NewFormItem("Search Pattern", searchPatternEntry)

	line := canvas.NewLine(color.White)
	form := widget.NewForm(radioFormItem, searchPathFormItem, searchPatternFormItem)
	form.SubmitText = "Search"
	form.CancelText = "Clear All"
	form.OnSubmit = formSubmitHandler
	form.OnCancel = formCancelHandler

	return container.New(layout.NewVBoxLayout(), line, form)
}
