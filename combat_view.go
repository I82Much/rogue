package rogue

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
)

type CombatView struct {
	Model *CombatModel
}

// Render assumes that termbox has already been initialized.
func (v *CombatView) Render() {
	// Draw all of the falling words
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, word := range v.Model.Words() {
		foreground := termbox.ColorDefault
		if word == v.Model.CurrentlyTyping() {
			foreground = foreground | termbox.AttrBold
		}
		for i, c := range word.word {
			if i < len(word.spelled) {
				termbox.SetCell(i, word.row, c, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
			} else {
				termbox.SetCell(i, word.row, c, foreground, termbox.ColorDefault)
			}
		}
	}

	// Draw the % accuracy
	if v.Model.attempts > 0 {
		hits := v.Model.hits
		attempts := v.Model.attempts
		accuracyText := fmt.Sprintf("%d / %d (%.2f%%)", hits, attempts, 100.0*float32(hits)/float32(attempts))
		for i, c := range accuracyText {
			termbox.SetCell(i, 3, c, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}
