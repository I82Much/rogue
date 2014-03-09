package rogue

import (
	"fmt"
)

type Tile int32

type Creature int32

type MovementResult int32

const (
	// Tiles
	Floor Tile = iota
	Wall

	// Creatures
	None Creature = iota
	Player
	Monster

	// Movement possibilities
	Move MovementResult = iota
	OutOfBounds
	Impassable
	CreatureOccupying
)

var (
	InvalidPlayerLoc = Location{Row: -1, Col: -1}
)

type World struct {
	rows, cols int
	tiles      [][]Tile
	creatures  [][]Creature
	playerLoc  Location
}

func NewWorld(rows, cols int) *World {
	tiles := make([][]Tile, rows)
	for row := 0; row < rows; row++ {
		tiles[row] = make([]Tile, cols)
	}
	creatures := make([][]Creature, rows)
	for row := 0; row < rows; row++ {
		creatures[row] = make([]Creature, cols)
		for col := 0; col < cols; col++ {
			creatures[row][col] = None
		}
	}
	return &World{
		rows:      rows,
		cols:      cols,
		tiles:     tiles,
		creatures: creatures,
		playerLoc: InvalidPlayerLoc,
	}
}

func (w *World) TileAt(loc Location) Tile {
	return w.tiles[loc.Row][loc.Col]
}

func (w *World) SetTile(loc Location, t Tile) {
	w.tiles[loc.Row][loc.Col] = t
}

func (w *World) CreatureAt(loc Location) Creature {
	return w.creatures[loc.Row][loc.Col]
}

func (w *World) SetCreature(loc Location, c Creature) MovementResult {
	res := w.CanMoveTo(loc)
	if res == Move {
		w.creatures[loc.Row][loc.Col] = c
	}
	return res
}

func (w *World) Rows() int {
	return w.rows
}

func (w *World) Cols() int {
	return w.cols
}

func (w *World) Spawn(row, col int) {
	if w.playerLoc != InvalidPlayerLoc {
		panic("player already spawned")
	}
	if w.SetCreature(Loc(row, col), Player) != Move {
		panic("player can't spawn here")
	}
	w.playerLoc = Loc(row, col)
}

// Attempts to spawn a monster
func (w *World) SpawnMonster() bool {
	for row := 0; row < w.Rows(); row++ {
		for col := 0; col < w.Cols(); col++ {
			if w.SetCreature(Loc(row, col), Monster) == Move {
				return true
			}
		}
	}
	return false
}

// MovePlayer moves the player the given number of rows/cols relative
// to where he already is. No-op if out of bounds / can't move there.
func (w *World) MovePlayer(rows, cols int) MovementResult {
	newLoc := w.playerLoc.Add(Location{Row: rows, Col: cols})
	res := w.SetCreature(newLoc, Player)
	if res == Move {
		// Remove old value
		w.creatures[w.playerLoc.Row][w.playerLoc.Col] = None
		w.playerLoc = newLoc
	}
	return res
}

func (w *World) InBounds(loc Location) bool {
	return loc.Row >= 0 && loc.Row < w.Rows() &&
		loc.Col >= 0 && loc.Col < w.Cols()
}

func (w *World) CanMoveTo(loc Location) MovementResult {
	// In bounds
	inBounds := w.InBounds(loc)
	if !inBounds {
		return OutOfBounds
	}

	// Is there a creature in that spot
	// TODO ndunn this probably needs to change for combat to work
	if got := w.CreatureAt(loc); got != None {
		return CreatureOccupying
	}

	passable := w.TileAt(loc).Passable()
	if !passable {
		return Impassable
	}
	return Move
}

func (t Tile) Rune() rune {
	switch t {
	case Floor:
		return ' '
	case Wall:
		return '*'
	default:
		panic(fmt.Sprintf("unknown tile type %v", t))
	}
}

func (t Tile) Passable() bool {
	if t == Floor {
		return true
	}
	return false
}

func (c Creature) Rune() rune {
	switch c {
	case None:
		return ' '
	case Player:
		return 'P'
	case Monster:
		return 'M'
	default:
		panic(fmt.Sprintf("unknown monster type %v", c))
	}
}

func (m MovementResult) String() string {
	switch m {
	case Move:
		return "Move"
	case OutOfBounds:
		return "Out of bounds"
	case Impassable:
		return "Impassable"
	case CreatureOccupying:
		return "Creature"
	default:
		panic(fmt.Sprintf("unknown movement result type %v", m))

	}
}

/*
// TODO(ndunn): rendering shouldn't be in world.
func (w World) String() string {
	rows := make([]string, w.Rows())
	for _, row := range w.tiles {
		rowString := ""
		for _, tile := range row {
			rowString += tile.String()
		}
		rows = append(rows, rowString)
	}
	return strings.Join(rows, "\n")
}*/

func (w *World) RuneAt(loc Location) rune {
	if c := w.CreatureAt(loc); c != None {
		return c.Rune()
	}
	return w.TileAt(loc).Rune()
}
