package main

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createTopContainer() *fyne.Container {

	line := canvas.NewLine(color.White)
	label := widget.NewLabel("Hello")
	return container.New(layout.NewVBoxLayout(), label, line)
}

func createMiddleContainer() {}

func main() {

	application := app.New()
	window := application.NewWindow("friendly-grep")

	data := binding.BindStringList(
		&[]string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5", "Item 6", "Item 6", "Item 8", "Item 9", "Item 10", "Item 11", "Item 12", "Item 13", "Item 14", "Item 15", "Item 16"},
	)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	sourceEntry := widget.NewEntry()
	searchEntry := widget.NewEntry()

	radio := widget.NewRadioGroup([]string{"File", "Directory"}, func(selected string) {
		log.Println(selected)
	})

	radio.SetSelected("File")

	browseFile := func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			log.Println(reader.URI().String())

		}, window)
	}

	browseFolder := func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if list == nil {
				log.Println("Cancelled")
				return
			}

			children, err := list.List()
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
			dialog.ShowInformation("Folder Open", out, window)
		}, window)
	}

	button := widget.NewButton("Browse", func() {
		if radio.Selected == "File" {
			browseFile()
		}

		if radio.Selected == "Directory" {
			browseFolder()
		}
	})

	h := container.New(layout.NewHBoxLayout(), button, radio)
	formItem0 := widget.NewFormItem("", h)
	formItem1 := widget.NewFormItem("", sourceEntry)
	formItem2 := widget.NewFormItem("Search Pattern", searchEntry)

	form := widget.NewForm(formItem0, formItem1, formItem2)

	// optional on-submit button
	form.OnSubmit = func() {
		log.Println("Form submitted:", sourceEntry.Text)
		log.Println("Form submitted:", searchEntry.Text)
		window.Close()
	}

	form.SubmitText = "Search"

	line2 := canvas.NewLine(color.White)
	top := createTopContainer()

	middle := list
	bottom := container.New(layout.NewVBoxLayout(), line2, form)

	content := container.New(layout.NewBorderLayout(top, bottom, nil, nil),
		top, bottom, middle)

	window.Resize(fyne.NewSize(900, 700))

	window.SetContent(content)
	window.Show()
	application.Run()

}
