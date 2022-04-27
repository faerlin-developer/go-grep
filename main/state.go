package main

import (
	"fmt"
	"go-grep/grep"
	"image/color"
	"os"
	"strconv"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type State struct {
	UserInput UserInput
	View      View
}

type UserInput struct {
	SelectedOption RadioOption
	SearchPath     *binding.String
	SearchPattern  *binding.String
	NumberWorkers  *binding.String
	BufferSize     *binding.String
}

type View struct {
	Scroll   *container.Scroll
	TextGrid *widget.TextGrid
	nextRow  int
}

const (
	File RadioOption = iota
	Directory
	Unrecognized
)

type RadioOption int

// ----- RadioOption -----

func (r RadioOption) Value() string {
	return [...]string{"File", "Directory", ""}[r]
}

// ----- View -----

func (v *View) SetText(row int, text string) {

	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		v.TextGrid.SetRune(row, i, runes[i])
	}
}

func (v *View) Clear() {
	textGrid := widget.NewTextGrid()
	textGrid.ShowLineNumbers = false
	v.TextGrid = textGrid
	v.Scroll.Content = textGrid
	v.nextRow = 0
	v.Scroll.Refresh()
}

func (v *View) ShowDefaultView() {
	state.View.Clear()
	cwd, _ := os.Getwd()
	state.View.AppendText("")
	state.View.AppendText("")
	state.View.AppendText("┌─┐┌─┐   ┌─┐┬─┐┌─┐┌─┐")
	state.View.AppendText("│ ┬│ │───│ ┬├┬┘├┤ ├─┘")
	state.View.AppendText("└─┘└─┘   └─┘┴└─└─┘┴   is a blah blah blah ")
	state.View.AppendText("")
	state.View.AppendText(fmt.Sprintf("Current working directory: %v", cwd))
	state.View.TextGrid.ShowLineNumbers = false
}

// a command-line utility for searching plain-text data sets for lines that match a regular expression.

// ----- UserInput -----

func (u *UserInput) SetSelectedOption(selectedOption string) {
	u.SelectedOption = toRadioOption(selectedOption)
}

func (u *UserInput) GetSelectedOption() RadioOption {
	return u.SelectedOption
}

func (u *UserInput) SetSearchPath(searchPath string) {
	(*(u.SearchPath)).Set(searchPath)
}

func (u *UserInput) GetSearchPath() string {
	value, _ := (*(u.SearchPath)).Get()
	return value
}

func (u *UserInput) ClearSearchPath() {
	(*(u.SearchPath)).Set("")
}

func (u *UserInput) GetSearchTerm() string {
	value, _ := (*(u.SearchPattern)).Get()
	return value
}

func (u *UserInput) ClearSearchTerm() {
	(*(u.SearchPattern)).Set("")
}

func (u *UserInput) SetNumberWorkers(numberWorkers int) {
	(*(u.NumberWorkers)).Set(strconv.Itoa(numberWorkers))
}

func (u *UserInput) GetNumberWorkers() int {
	value, _ := (*(u.NumberWorkers)).Get()
	intValue, _ := strconv.Atoi(value)
	return intValue
}

func (u *UserInput) SetBufferSize(bufferSize int) {
	(*(u.BufferSize)).Set(strconv.Itoa(bufferSize))
}

func (u *UserInput) GetBufferSize() int {
	value, _ := (*(u.BufferSize)).Get()
	intValue, _ := strconv.Atoi(value)
	return intValue
}

func (v *View) AppendText(text string) {
	nextRow := v.nextRow
	v.SetText(nextRow, text)
	v.nextRow += 1
}

func (v *View) AppendResult(filepath string, line string, lineNumber int, indices [][]int) {

	result := fmt.Sprintf("%v:%v:%v", filepath, lineNumber, line)
	v.AppendText(result)

	v.TextGrid.SetStyleRange(v.nextRow-1, 0, v.nextRow-1, len(filepath)-1,
		&widget.CustomTextGridStyle{FGColor: &color.NRGBA{R: 51, G: 153, B: 255, A: 255}})

	numDigit := getNumDigits(lineNumber)

	for _, idx := range indices {
		v.TextGrid.SetStyleRange(v.nextRow-1, len(filepath)+1, v.nextRow-1,
			len(filepath)+1+numDigit, &widget.CustomTextGridStyle{FGColor: &color.NRGBA{R: 0, G: 255, B: 0, A: 255}})
		v.TextGrid.SetStyleRange(v.nextRow-1, len(filepath)+2+numDigit+idx[0], v.nextRow-1,
			len(filepath)+2+numDigit+idx[1]-1, &widget.CustomTextGridStyle{FGColor: &color.NRGBA{R: 255, G: 204, B: 0, A: 255}})
	}
}

// ----- Public Utility Functions ----

func NewState() State {
	userInput := NewUserInput()
	view := NewView()
	userInput.SetNumberWorkers(grep.DefaultNumberWorkers)
	userInput.SetBufferSize(grep.DefaultBufferSize)
	return State{UserInput: userInput, View: view}
}

func NewView() View {
	textGrid := widget.NewTextGrid()
	textGrid.ShowLineNumbers = false
	scroll := container.NewScroll(textGrid)
	return View{Scroll: scroll, TextGrid: textGrid, nextRow: 0}
}

func NewUserInput() UserInput {
	searchPath := binding.NewString()
	searchPattern := binding.NewString()
	numberWorkers := binding.NewString()
	bufferSize := binding.NewString()
	return UserInput{SelectedOption: File, SearchPath: &searchPath,
		SearchPattern: &searchPattern, NumberWorkers: &numberWorkers, BufferSize: &bufferSize}
}

// ----- Private Utility Functions -----

func getNumDigits(x int) int {
	x_as_string := strconv.Itoa(x)
	return len(x_as_string)
}

func toRadioOption(value string) RadioOption {

	var radioOption RadioOption
	switch value {
	case File.Value():
		radioOption = File
	case Directory.Value():
		radioOption = Directory
	default:
		radioOption = Unrecognized
	}

	return radioOption
}
