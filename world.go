package rogue

import (
	"fmt"
	"strings"
)

type Tile int32

// TODO(ndunn): player isn't really a tile.
const (
	FLOOR Tile = iota
	WALL
	PLAYER
)

type World struct {
	rows, cols int
	tiles      [][]Tile
}

func NewWorld(rows, cols int) World {
	tiles := make([][]Tile, rows)
	for row := 0; row < rows; row++ {
		tiles[row] = make([]Tile, cols)
	}
	return World{
		rows:  rows,
		cols:  cols,
		tiles: tiles,
	}
}

func (w World) At(row, col int) Tile {
	return w.tiles[row][col]
}

func (w World) Set(row, col int, t Tile) {
	w.tiles[row][col] = t
}

func (w World) Rows() int {
	return w.rows
}

func (w World) Cols() int {
	return w.cols
}

func (t Tile) String() string {
	switch t {
	case FLOOR:
		return " "
	case WALL:
		return "*"
	case PLAYER:
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
