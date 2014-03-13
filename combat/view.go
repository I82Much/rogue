package combat

import (
	"fmt"
	"strings"

	"github.com/I82Much/rogue/math"
	termbox "github.com/nsf/termbox-go"
)

type View struct {
	Model *Model
	rows  int
}

// Render assumes that termbox has already been initialized.
func (v *View) Render() {
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
		healthWidth := math.IntMap(monster.Life, 0, monster.MaxLife, 0, healthBarWidth)
		if healthWidth < 0 {
			healthWidth = 0
		}
		for h := 0; h < healthWidth; h++ {
			termbox.SetCell(offset+5+h, 0, '█', termbox.ColorRed, termbox.ColorDefault)
		}
	}

	// Render the player

	// TODO(ndunn): all of the monster graphics etc should be moved into template files
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
	// Draw player's health bar
	healthWidth := 40
	player := v.Model.Player
	life := math.IntMap(player.Life, 0, player.MaxLife, 0, healthWidth)
	if life < 0 {
		life = 0
	}
	for j := 0; j < healthWidth; j++ {
		row := 15
		if j < life {
			termbox.SetCell(j, row, '█', termbox.ColorRed, termbox.ColorDefault)
		} else {
			termbox.SetCell(j, row, '░', termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	for _, word := range v.Model.Words() {
		foreground := termbox.ColorDefault
		if word == v.Model.CurrentlyTyping() {
			foreground = foreground | termbox.AttrBold
		}

		row := 0
		// TODO(ndunn): render some sort of line to show where the dividing point is
		numRows := 25
		// If we're attacking, words are flying up towards the enemies
		if v.Model.State() == Attack {
			row = int(math.Lerp(float64(numRows), 0.0, word.proportion))
		} else if v.Model.State() == Defense {
			row = int(math.Lerp(0.0, float64(numRows), word.proportion))
		}

		// If we're defending, words are flying down towards player
		for i, c := range word.word {
			if i < len(word.spelled) {
				termbox.SetCell(i, row, c, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
			} else {
				termbox.SetCell(i, row, c, foreground, termbox.ColorDefault)
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
