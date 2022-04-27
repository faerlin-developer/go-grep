package main

import (
	"go-grep/grep"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createMenu() *fyne.MainMenu {

	settingsMenuItem := fyne.NewMenuItem("Settings", settingsMenuItemHandler)
	menu := fyne.NewMenu("Menu", settingsMenuItem)
	return fyne.NewMainMenu(menu)
}

func createMiddleComponent() *container.Scroll {
	state.View.ShowDefaultView()
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

func settingsMenuItemHandler() {
	settingsWindow := application.NewWindow("Settings")

	numberWorkersEntry := NewNumericalEntry(state.UserInput.NumberWorkers)
	bufferSizeEntry := NewNumericalEntry(state.UserInput.BufferSize)

	numberFormItem := widget.NewFormItem("Number of Workers", numberWorkersEntry)
	bufferSizeFormItem := widget.NewFormItem("Buffer Size", bufferSizeEntry)
	form := widget.NewForm(numberFormItem, bufferSizeFormItem)

	form.OnSubmit = func() {
		state.UserInput.SetNumberWorkers(grep.DefaultNumberWorkers)
		state.UserInput.SetBufferSize(grep.DefaultBufferSize)
	}

	form.SubmitText = "Set Default Values"

	settingsWindow.SetContent(form)
	settingsWindow.CenterOnScreen()
	settingsWindow.Show()
}
