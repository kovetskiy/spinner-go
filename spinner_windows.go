package spinner

import termbox "github.com/nsf/termbox-go"

func getTerminalWidth() int {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	width, _ := termbox.Size()
	termbox.Close()
	return width
}
