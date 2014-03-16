package dungeon

import (
	"fmt"
	"strconv"

	"github.com/I82Much/rogue/monster"
	"github.com/I82Much/rogue/render"

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

func (t Tile) Passable() bool {
	switch t {
	case Floor, UnlockedDoor, Bridge:
		return true
	}
	return false
}

func (t Tile) Rune() rune {
	switch t {
	case Floor:
		return ' '
	case Wall:
		return '*'
	case UnlockedDoor:
		return 'D'
	case LockedDoor:
		return 'L'
	case Water:
		return '░'
	case Bridge:
		return '█'
	default:
		panic(fmt.Sprintf("unknown tile type %v", t))
	}
}

func (t Tile) Foreground() termbox.Attribute {
	switch t {
	case Water:
		return termbox.ColorWhite
	case Bridge:
		return termbox.ColorYellow
	}
	return termbox.ColorDefault
}

func (t Tile) Background() termbox.Attribute {
	switch t {
	case Water:
		return termbox.ColorBlue
	}
	return termbox.ColorDefault
}

func RuneForMonsters(m []monster.Type) rune {
	// TODO(ndunn): It would be pretty cool if it used '?' for monsters you haven't fought yet
	if len(m) == 0 {
		return ' '
	} else if len(m) > 1 {
		return []rune(strconv.Itoa(len(m)))[0]
	}
	switch m[0] {
	case monster.Haxor, monster.HaxorScammer, monster.HaxorSpammer, monster.HaxorBlogger:
		return 'H'
	case monster.Scammer:
		return '$'
	case monster.Spammer:
		return 'S'
	case monster.Blogger:
		return 'B'
	}
	panic(fmt.Sprintf("unknown type %v", m[0]))
}

type cell struct {
	r  rune
	fg termbox.Attribute
	bg termbox.Attribute
}

func (w *Room) CellForLoc(loc Location) cell {
	if w.playerLoc == loc {
		return cell{
			r:  'P',
			fg: termbox.ColorDefault,
			bg: termbox.ColorDefault,
		}
	}
	if m := w.MonstersAt(loc); m != nil {
		return cell{
			r:  RuneForMonsters(m),
			fg: termbox.ColorRed,
			bg: termbox.ColorDefault,
		}
	}
	t := w.TileAt(loc)
	return cell{
		r:  t.Rune(),
		fg: t.Foreground(),
		bg: t.Background(),
	}
}

func (v *View) renderRoom(r *Room) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for row := 0; row < r.Rows(); row++ {
		for col := 0; col < r.Cols(); col++ {
			// col = x, row = y
			location := Loc(row, col)
			cell := r.CellForLoc(location)
			termbox.SetCell(col, row, cell.r, cell.fg, cell.bg)
		}
	}

	// Render data about the player
	player := v.model.player
	render.Render(player.Name, r.Rows(), 0)
	// Health
	render.Render(fmt.Sprintf("%d/%d health", player.CurrentLife, player.MaxLife), r.Rows()+1, 0)
	render.Render(fmt.Sprintf("%v", player.Stats), r.Rows()+2, 0)
	termbox.Flush()
}

func (v *View) Render() {
	v.renderRoom(v.model.world.CurrentRoom())
}
