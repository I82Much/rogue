package render

import (
	"strings"

	termbox "github.com/nsf/termbox-go"
)

func RenderWithColor(msg string, startRow, startCol int, foreground termbox.Attribute, background termbox.Attribute) {
	for row, line := range strings.Split(msg, "\n") {
		for i, char := range line {
			termbox.SetCell(startCol+i, startRow+row, char, foreground, background)
		}
	}
}

func Render(msg string, startRow, startCol int) {
	RenderWithColor(msg, startRow, startCol, termbox.ColorDefault, termbox.ColorDefault)
}
