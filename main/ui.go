package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NewWindow() fyne.Window {
	return app.New().NewWindow("go-grep")
}

func createTopComponent() *fyne.Container {

	line := canvas.NewLine(color.White)
	label := widget.NewLabel("Hello")
	return container.New(layout.NewVBoxLayout(), label, line)
}

func createMiddleComponent() *widget.List {

	createItem := func() fyne.CanvasObject {
		return widget.NewLabel("template")
	}

	updateItem := func(i binding.DataItem, o fyne.CanvasObject) {
		o.(*widget.Label).Bind(i.(binding.String))
	}

	return widget.NewListWithData(state.Data, createItem, updateItem)
}

func createBottomComponent() *fyne.Container {

	button := widget.NewButton("Browse", browseButtonHandler)
	radio := widget.NewRadioGroup([]string{File.Value(), Directory.Value()}, radioGroupHandler)
	radio.SetSelected(File.Value())

	searchPathEntry := widget.NewEntryWithData(*state.SearchPath)
	searchTermEntry := widget.NewEntryWithData(*state.SearchTerm)
	browseContainer := container.New(layout.NewHBoxLayout(), button, radio)

	browseFormItem := widget.NewFormItem("", browseContainer)
	searchPathFormItem := widget.NewFormItem("Search Path", searchPathEntry)
	searchTermFormItem := widget.NewFormItem("Search Pattern", searchTermEntry)

	line := canvas.NewLine(color.White)
	form := widget.NewForm(browseFormItem, searchPathFormItem, searchTermFormItem)
	form.SubmitText = "Search"
	form.CancelText = "Clear All"
	form.OnSubmit = formSubmitHandler
	form.OnCancel = formCancelHandler

	return container.New(layout.NewVBoxLayout(), line, form)
}
