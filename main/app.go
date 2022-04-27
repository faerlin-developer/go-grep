package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

var state State
var application fyne.App
var window fyne.Window

func init() {
	state = NewState()
	application = app.New()
	window = application.NewWindow("go-grep")
}

func main() {

	menu := createMenu()
	middle := createMiddleComponent()
	bottom := createBottomComponent()
	content := container.New(layout.NewBorderLayout(nil, bottom, nil, nil),
		bottom, middle)

	application.Settings().SetTheme(theme.DarkTheme())

	window.SetMainMenu(menu)
	window.SetMaster()
	window.CenterOnScreen()
	window.SetContent(content)
	window.Resize(fyne.NewSize(900, 700))
	window.ShowAndRun()
}
