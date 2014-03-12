package render

import (
	"strings" 
	
	termbox "github.com/nsf/termbox-go"
)

func Render(msg string, startRow, startCol int) {
	for row, line := range strings.Split(msg, "\n") {
		for i, char := range line {
			termbox.SetCell(startCol + i, startRow + row, char, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}