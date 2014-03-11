package rogue

import (
	"fmt"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

type CombatView struct {
	Model *CombatModel
	rows  int
}

// Render assumes that termbox has already been initialized.
func (v *CombatView) Render() {
	// Draw all of the falling words
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Render the monsters
	for i, monster := range v.Model.Monsters {
		monsterFigure := `
		 \ O /
		   |
		  / \`
		healthBarWidth := 25

		offset := i * (10 + healthBarWidth)

		for row, figure := range strings.Split(monsterFigure, "\n") {
			for j, char := range figure {
				termbox.SetCell(offset+j, row, char, termbox.ColorDefault, termbox.ColorDefault)
			}
		}
		// Draw the health bar
		healthWidth := intMap(monster.Life, 0, monster.MaxLife, 0, healthBarWidth)
		if healthWidth < 0 {
			healthWidth = 0
		}
		for h := 0; h < healthWidth; h++ {
			termbox.SetCell(offset+5+h, 0, '█', termbox.ColorRed, termbox.ColorDefault)
		}
	}

	// Render the player

	playerFigure := `
                                               +
                +--------------+               |
                |              |               |
                |              |               |
+---------+     |              |               |
|         |     |              |               |
|         |     |              |            +----->
|         |     |              |               |
|         | +----------------------+           |
|    ^    | |                      |           |
|    |    | |                      |           |
|    |    | |                      +-----------+
+----|----+ |                      +
     +------+                      |
            |                      |
            |                      |
            +----------------------+
	
	`
	for row, figure := range strings.Split(playerFigure, "\n") {
		for j, char := range figure {
			finalRow := 10 + row
			termbox.SetCell(j, finalRow, char, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

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
