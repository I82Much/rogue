package combat

import (
	"fmt"
	"log"
	"strings"

	"github.com/I82Much/rogue/math"
	"github.com/I82Much/rogue/monster"
	"github.com/I82Much/rogue/render"
	"github.com/I82Much/rogue/static"
	termbox "github.com/nsf/termbox-go"
)

type View struct {
	Model       *Model
	rows        int
	description renderer
}

type renderer interface {
	Render()
}

func NewView(m *Model, rows int) *View {
	return &View{
		Model:       m,
		rows:        rows,
		description: initModule(m),
	}
}

func columnOffset(w *AttackWord) int {
	// TODO(ndunn): this really needs to know the width of the window, the length of word, etc.
	switch w.Col {
	case Left:
		return 5
	case Center:
		return 20
	case Right:
		return 35
	default:
		panic("unknown column")
	}
}

func initModule(m *Model) renderer {
	var monsterMap = make(map[monster.Type]int)
	for _, monster := range m.Monsters {
		monsterMap[monster.Type]++
	}
	text := "Get ready to fight "
	var descriptions []string
	for t, num := range monsterMap {
		descriptions = append(descriptions, t.Description(num))
	}
	text += strings.Join(descriptions, ",")
	// TODO would be good to be able to skip..
	return static.NewModule(text, nil)
}

func (v *View) RenderInitial() {
	v.description.Render()
}


func (v *View) RenderWords() {
	for _, word := range v.Model.Words() {
		// Some words aren't actually visible yet - they're in the model but there's a delay
		if !word.IsVisible() {
			log.Printf("word %v is not visible", word.word)
			continue
		}

		foreground := termbox.ColorDefault
		if word == v.Model.CurrentlyTyping() {
			foreground = foreground | termbox.AttrBold
		}

		row := 0
		// TODO(ndunn): render some sort of line to show where the dividing point is
		// If we're attacking, words are flying up towards the enemies
		if v.Model.State() == Attack {
			row = int(math.Lerp(float64(v.rows), 0.0, word.proportion))
		} else if v.Model.State() == Defense {
			row = int(math.Lerp(0.0, float64(v.rows), word.proportion))
		}

		colOffset := columnOffset(word)
		// If we're defending, words are flying down towards player
		for i, c := range word.word {
			col := colOffset + i
			if i < len(word.spelled) {
				termbox.SetCell(col, row, c, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
			} else {
				termbox.SetCell(col, row, c, foreground, termbox.ColorDefault)
			}
		}
	}
}

func (v *View) RenderMonsters() {
	// Render the monsters
	for i, monster := range v.Model.Monsters {
		monsterFigure := `
		 \ O /
		   |
		  / \`
		healthBarWidth := 25

		offset := i * (10 + healthBarWidth)
		row := 0
		render.Render(monsterFigure, row, offset)
		/*for row, figure := range strings.Split(monsterFigure, "\n") {
			for j, char := range figure {
				termbox.SetCell(offset+j, row, char, termbox.ColorDefault, termbox.ColorDefault)
			}
		}*/
		// Draw the health bar
		healthWidth := math.IntMap(monster.Life, 0, monster.MaxLife, 0, healthBarWidth)
		if healthWidth < 0 {
			healthWidth = 0
		}
		for h := 0; h < healthWidth; h++ {
			termbox.SetCell(offset+5+h, 0, '█', termbox.ColorRed, termbox.ColorDefault)
		}
	}
}

func (v *View) RenderPlayer() {
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
	life := math.IntMap(player.CurrentLife, 0, player.MaxLife, 0, healthWidth)
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
}

func (v *View) RenderAccuracy() {
	// Draw the % accuracy
	if v.Model.attempts > 0 {
		hits := v.Model.hits
		attempts := v.Model.attempts
		accuracyText := fmt.Sprintf("%d / %d (%.2f%%)", hits, attempts, 100.0*float32(hits)/float32(attempts))
		for i, c := range accuracyText {
			termbox.SetCell(50+i, 30, c, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

func (v *View) RenderCombat() {
	// Draw all of the falling words
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	v.RenderMonsters()
	v.RenderPlayer()

	// Draw the dividing line, below which the word does damage
	divider := "________________________________________________________________"
	// Pull it up one row so that it is at the TOP of where it can be.
	render.RenderWithColor(divider, v.rows - 1, 0, termbox.ColorRed, termbox.ColorDefault)

	// Falling/rising words
	v.RenderWords()
	v.RenderAccuracy()
	termbox.Flush()
}

// Render assumes that termbox has already been initialized.
func (v *View) Render() {
	if v.Model.State() == EnemyDescription {
		v.RenderInitial()
	} else {
		v.RenderCombat()
	}

}
