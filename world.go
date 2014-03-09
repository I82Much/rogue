package rogue

type Tile int32

const (
	FLOOR Tile = iota
	WALL
)

type World struct {
	tiles [][]Tile
}

func NewWorld(rows, cols int) World {
	tiles := make([][]Tile, rows)
	for row := 0; row < rows; row++ {
		tiles[row] = make([]Tile, cols)
	}
	return World{
		tiles: tiles,
	}
}
