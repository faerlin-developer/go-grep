package main

import (
	"fyne.io/fyne/v2/data/binding"
)

type RadioOption int

const (
	File RadioOption = iota
	Directory
	Unrecognized
)

type State struct {
	SelectedOption RadioOption
	SearchPath     *binding.String
	SearchTerm     *binding.String
	Data           binding.ExternalStringList
}

func (r RadioOption) Value() string {
	return [...]string{"File", "Directory", ""}[r]
}

func (s *State) SetSelectedOption(selectedOption string) {
	s.SelectedOption = toRadioOption(selectedOption)
}

func (s *State) GetSelectedOption() RadioOption {
	return s.SelectedOption
}

func (s *State) SetSearchPath(searchPath string) {
	(*(s.SearchPath)).Set(searchPath)
}

func (s *State) GetSearchPath() string {
	value, _ := (*(s.SearchPath)).Get()
	return value
}

func (s *State) ClearSearchPath() {
	(*(s.SearchPath)).Set("")
}

func (s *State) SetSearchTerm(searchTerm string) {
	(*(s.SearchTerm)).Set(searchTerm)
}

func (s *State) GetSearchTerm() string {
	value, _ := (*(s.SearchTerm)).Get()
	return value
}

func (s *State) ClearSearchTerm() {
	(*(s.SearchTerm)).Set("")
}

func (s *State) ClearData() {
	s.Data.Set([]string{})
}

func (s *State) AppendData(line string) {
	s.Data.Append(line)
}

func NewState() State {
	searchPath := binding.NewString()
	searchTerm := binding.NewString()
	data := binding.BindStringList(&[]string{"Item 1", "Item 2", "Item 3"})
	return State{SelectedOption: File, SearchPath: &searchPath, SearchTerm: &searchTerm, Data: data}
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
