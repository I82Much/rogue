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
