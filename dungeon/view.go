package dungeon

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
)

type View struct {
	model *Model
}

func NewView(m *Model) *View {
	return &View{
		model: m,
	}
}

func (t Tile) Rune() rune {
	switch t {
	case Floor:
		return ' '
	case Wall:
		return '*'
	case DoorTile:
		return 'D'
	default:
		panic(fmt.Sprintf("unknown tile type %v", t))
	}
}

func (c Creature) Rune() rune {
	switch c {
	case None:
		return ' '
	case PlayerCreature:
		return 'P'
	case MonsterCreature:
		return 'M'
	default:
		panic(fmt.Sprintf("unknown monster type %v", c))
	}
}

func (w *Room) RuneAt(loc Location) rune {
	if c := w.CreatureAt(loc); c != None {
		return c.Rune()
	}
	return w.TileAt(loc).Rune()
}

func renderRoom(r *Room) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for row := 0; row < r.Rows(); row++ {
		for col := 0; col < r.Cols(); col++ {
			// col = x, row = y
			location := Loc(row, col)
			termbox.SetCell(col, row, r.RuneAt(location), termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func (v *View) Render() {

	renderRoom(v.model.world.CurrentRoom())
}
