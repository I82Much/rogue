package dungeon

import (
	"fmt"

	"github.com/I82Much/rogue/monster"
)

type Tile int32

type Creature int32

type MovementResult int32

const (
	// Tiles
	Floor Tile = iota
	Wall
	LockedDoor
	UnlockedDoor
	Water
	Bridge

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
	//creatures  [][]Creature

	monsters  map[Location][]monster.Type
	playerLoc Location
	// Sparse map
	doors map[Location]*Door
}

func NewRoom(rows, cols int) *Room {
	tiles := make([][]Tile, rows)
	for row := 0; row < rows; row++ {
		tiles[row] = make([]Tile, cols)
	}
	return &Room{
		rows:      rows,
		cols:      cols,
		tiles:     tiles,
		monsters:  make(map[Location][]monster.Type),
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
	w.SetTile(loc, LockedDoor)
}

func (w *Room) MonstersAt(loc Location) []monster.Type {
	return w.monsters[loc]
}

func (w *Room) RemovePlayer() {
	if w.playerLoc == InvalidLoc {
		return
	}
	w.playerLoc = InvalidLoc
}

// After combat, we take the place where the monster was formerly occupying
func (w *Room) ReplaceMonsterWithPlayer(loc Location) {
	delete(w.monsters, loc)
	/*w.creatures[loc.Row][loc.Col] = PlayerCreature
	w.creatures[w.playerLoc.Row][w.playerLoc.Col] = None*/
	w.playerLoc = loc
}

func (w *Room) UnlockDoors() {
	for loc := range w.doors {
		w.SetTile(loc, UnlockedDoor)
	}
}

func (w *Room) AnyMonstersLeft() bool {
	if w == nil {
		return false
	}
	return len(w.monsters) > 0
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
	if reason := w.PlayerCanOccupy(Loc(row, col)); reason != Move {
		panic(fmt.Sprintf("player can't spawn here: %v", reason))
	}
	w.playerLoc = Loc(row, col)
}

// Returns error if it couldn't add monster
func (w *Room) AddMonster(row, col int, m monster.Type) error {
	loc := Loc(row, col)
	if res := w.MonsterCanOccupy(loc); res != Move {
		return fmt.Errorf("Monster can't occupy %v: %v", loc, res)
	}
	w.monsters[loc] = append(w.monsters[loc], m)
	return nil
}

// MovePlayer moves the player the given number of rows/cols relative
// to where he already is. No-op if out of bounds / can't move there.
func (w *Room) MovePlayer(rows, cols int) MovementResult {
	newLoc := w.playerLoc.Add(Location{Row: rows, Col: cols})
	res := w.PlayerCanOccupy(newLoc)
	if res == Move {
		w.playerLoc = newLoc
	}
	return res
}

func (w *Room) PlayerCanOccupy(loc Location) MovementResult {
	// Players can occupy all the same spaces as monsters, except they
	// can't move to a place that a monster exists (without starting a fight)
	if w.MonstersAt(loc) != nil {
		return CreatureOccupying
	}
	return w.MonsterCanOccupy(loc)
}

func (w *Room) MonsterCanOccupy(loc Location) MovementResult {
	if !w.InBounds(loc) {
		return OutOfBounds
	}
	if !w.TileAt(loc).Passable() {
		return Impassable
	}
	return Move
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
	if got := w.MonstersAt(loc); got != nil {
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
