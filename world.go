package rogue

type Tile int32

const (
	FLOOR Tile = iota
	WALL
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
	return w[row][col]
}

func (w World) Rows() int {
	return w.rows
}

func (w World) Cols() int {
	return w.cols
}
