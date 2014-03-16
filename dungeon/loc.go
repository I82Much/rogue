package dungeon

type Location struct {
	Row, Col int
}

func (loc Location) Add(offset Location) Location {
	return Location{
		Row: loc.Row + offset.Row,
		Col: loc.Col + offset.Col,
	}
}

// Loc is convenience shorthand for constructing Location objects.
func Loc(rows, cols int) Location {
	return Location{
		Row: rows,
		Col: cols,
	}
}

func (l Location) ManhattanDist(loc Location) int {
	rowDiff := max(l.Row, loc.Row) - min(l.Row, loc.Row)
	colDiff := max(l.Col, loc.Col) - min(l.Col, loc.Col)
	return rowDiff + colDiff

}
