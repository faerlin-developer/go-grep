# __Go-grep__

__Go-grep__ is a graphical utility for searching text files for lines that match a provided regular expression. It is implemented in __Go__ and makes use of the [fyne](https://fyne.io/) module for its graphical interface. It takes a file or a directory as a search path and a regular expression as a search term. If the search path is a directory, file processing is performed by concurrent light-weight threads known as _goroutines_ communicating via _channels_. The default number of goroutines and buffer size of the channels are 10 and 100, respectively, and can be adjusted in the settings menu. 

<img src="img/main_window.png" alt="drawing" width="400"/><img src="img/result_window.png" alt="drawing" width="400"/>

## Built With

- [Go](https://go.dev/)
- [fyne](https://fyne.io/)