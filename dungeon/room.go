package dungeon

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
	DoorTile
	Water
	Bridge

	// Creatures
	None Creature = iota
	PlayerCreature
	MonsterCreature

	// Movement possibilities
	Move MovementResult = iota
	OutOfBounds
	Impassable
	CreatureOccupying
)

var (
	InvalidLoc = Location{Row: -1, Col: -1}
)

type Door struct {
	From, To *Room
	// This door in a given room is the same as this other door in another room. This allows the world
	// to place the player appropriately in the new room
	Same *Door
}

type Room struct {
	rows, cols int
	tiles      [][]Tile
	creatures  [][]Creature
	playerLoc  Location
	// Sparse map
	doors map[Location]*Door
}

func NewRoom(rows, cols int) *Room {
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
	return &Room{
		rows:      rows,
		cols:      cols,
		tiles:     tiles,
		creatures: creatures,
		playerLoc: InvalidLoc,
		doors:     make(map[Location]*Door),
	}
}

func WalledRoom(rows, cols int) *Room {
	r := NewRoom(rows, cols)
	// Make the outline set to WALL
	for i := 0; i < cols; i++ {
		// Top row
		r.SetTile(Loc(0, i), Wall)
		// Bottom row
		r.SetTile(Loc(rows-1, i), Wall)
	}
	for i := 0; i < rows; i++ {
		// Top row
		r.SetTile(Loc(i, 0), Wall)
		// Bottom row
		r.SetTile(Loc(i, cols-1), Wall)
	}
	return r
}

// DoorLocation returns the location of the given door, or nil if it doesn't exist
func (w *Room) DoorLocation(d *Door) *Location {
	for loc, door := range w.doors {
		if door == d {
			loc := loc
			return &loc
		}
	}
	return nil
}

func (w *Room) PlayerTile() Tile {
	return w.TileAt(w.playerLoc)
}

func (w *Room) TileAt(loc Location) Tile {
	return w.tiles[loc.Row][loc.Col]
}

func (w *Room) SetTile(loc Location, t Tile) {
	w.tiles[loc.Row][loc.Col] = t
}

func (w *Room) SetDoor(loc Location, d *Door) {
	w.doors[loc] = d
	w.SetTile(loc, DoorTile)
}

func (w *Room) CreatureAt(loc Location) Creature {
	return w.creatures[loc.Row][loc.Col]
}

func (w *Room) SetCreature(loc Location, c Creature) MovementResult {
	res := w.CanMoveTo(loc)
	if res == Move {
		w.creatures[loc.Row][loc.Col] = c
	}
	return res
}

func (w *Room) RemovePlayer() {
	if w.playerLoc == InvalidLoc {
		return
	}
	w.creatures[w.playerLoc.Row][w.playerLoc.Col] = None
	w.playerLoc = InvalidLoc
}

// After combat, we take the place where the monster was formerly occupying
func (w *Room) ReplaceMonsterWithPlayer(loc Location) {
	w.creatures[loc.Row][loc.Col] = PlayerCreature
	w.creatures[w.playerLoc.Row][w.playerLoc.Col] = None
	w.playerLoc = loc
}

func (w *Room) Rows() int {
	return w.rows
}

func (w *Room) Cols() int {
	return w.cols
}

func (w *Room) Spawn(row, col int) {
	if w.playerLoc != InvalidLoc {
		panic("player already spawned")
	}
	if w.SetCreature(Loc(row, col), PlayerCreature) != Move {
		panic("player can't spawn here")
	}
	w.playerLoc = Loc(row, col)
}

// Attempts to spawn a monster
func (w *Room) SpawnMonster() bool {
	for row := 0; row < w.Rows(); row++ {
		for col := 0; col < w.Cols(); col++ {
			if w.SetCreature(Loc(row, col), MonsterCreature) == Move {
				return true
			}
		}
	}
	return false
}

// MovePlayer moves the player the given number of rows/cols relative
// to where he already is. No-op if out of bounds / can't move there.
func (w *Room) MovePlayer(rows, cols int) MovementResult {
	newLoc := w.playerLoc.Add(Location{Row: rows, Col: cols})
	res := w.SetCreature(newLoc, PlayerCreature)
	if res == Move {
		// Remove old value
		w.creatures[w.playerLoc.Row][w.playerLoc.Col] = None
		w.playerLoc = newLoc
	}
	return res
}

func (w *Room) InBounds(loc Location) bool {
	return loc.Row >= 0 && loc.Row < w.Rows() &&
		loc.Col >= 0 && loc.Col < w.Cols()
}

func (w *Room) CanMoveTo(loc Location) MovementResult {
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
