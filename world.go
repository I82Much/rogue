package rogue

import (
	"fmt"
	"strings"
)

type Tile int32

// TODO(ndunn): player isn't really a tile.
const (
	Floor Tile = iota
	Wall
	Player
)

var (
	InvalidPlayerLoc = Location{Row: -1, Col: -1}
)

type World struct {
	rows, cols int
	tiles      [][]Tile
	playerLoc  Location
	oldTile    Tile
}

func NewWorld(rows, cols int) World {
	tiles := make([][]Tile, rows)
	for row := 0; row < rows; row++ {
		tiles[row] = make([]Tile, cols)
	}
	return World{
		rows:      rows,
		cols:      cols,
		tiles:     tiles,
		playerLoc: InvalidPlayerLoc,
	}
}

func (w World) At(loc Location) Tile {
	return w.tiles[loc.Row][loc.Col]
}

func (w World) Set(loc Location, t Tile) {
	w.tiles[loc.Row][loc.Col] = t
}

func (w World) Rows() int {
	return w.rows
}

func (w World) Cols() int {
	return w.cols
}

func (w World) Spawn(row, col int) {
	if w.playerLoc != InvalidPlayerLoc {
		panic("player already spawned")
	}
	w.oldTile = w.tiles[row][col]
	w.Set(Loc(row, col), Player)
	w.playerLoc = Loc(row, col)
}

// MovePlayer moves the player the given number of rows/cols relative
// to where he already is. No-op if out of bounds / can't move there.
func (w World) MovePlayer(rows, cols int) {
	newLoc := w.playerLoc.Add(Location{Row: rows, Col: cols})
	if !w.CanMoveTo(newLoc) {
		return
	}
	// Restore the old tile
	w.Set(w.playerLoc, w.oldTile)
	w.playerLoc = newLoc
	// Save old value of tile matrix and put the player there
	w.oldTile = w.At(newLoc)
	w.Set(newLoc, Player)
}

func (w World) CanMoveTo(loc Location) bool {
	// In bounds
	inBounds := loc.Row >= 0 && loc.Row < w.Rows() &&
		loc.Col >= 0 && loc.Col < w.Cols()
	passable := w.At(loc) == Floor
	return inBounds && passable
}

func (t Tile) String() string {
	switch t {
	case Floor:
		return " "
	case Wall:
		return "*"
	case Player:
		return "P"
	default:
		panic(fmt.Sprintf("unknown tile type %v", t))
	}
}

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
}
