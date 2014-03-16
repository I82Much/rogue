package dungeon

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/I82Much/rogue/monster"
)

type Tile int32

type Creature int32

type MovementResult int32

const (
	// Need 1 for each wall. But if we have monsters etc that's crowded.
	MinRows = 6
	MaxRows = 16
	MinCols = 6
	MaxCols = 40

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
	NextToDoor
)

var (
	InvalidLoc = Location{Row: -1, Col: -1}
)

type DoorDir string

const (
	North DoorDir = "NORTH"
	East  DoorDir = "EAST"
	South DoorDir = "SOUTH"
	West  DoorDir = "WEST"
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

func (r *Room) String() string {
	var rows []string
	for _, row := range r.tiles {
		tileRows := []string{}
		for _, col := range row {
			tileRows = append(tileRows, string(col.Rune()))
		}
		rows = append(rows, strings.Join(tileRows, ""))
	}
	return strings.Join(rows, "\n")
}

func NewRoom(rows, cols int) *Room {
	tiles := make([][]Tile, rows)
	for row := 0; row < rows; row++ {
		tiles[row] = make([]Tile, cols)
		for col := 0; col < cols; col++ {
			tiles[row][col] = Floor
		}
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

func (w *Room) SetDoor(dir DoorDir, d *Door) {
	loc := w.LocForDoor(dir)
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

func (w *Room) RandSpawn() {
	log.Printf("trying to find a space to spawn in\n %v", w.String())
	var open []Location
	for row := 0; row < w.Rows(); row++ {
		for col := 0; col < w.Cols(); col++ {
			loc := Loc(row, col)
			if w.PlayerCanOccupy(loc) == Move {
				open = append(open, loc)
			}
		}
	}
	if len(open) == 0 {
		panic("Cannot find any place to spawn")
	}
	loc := open[int(rand.Int31n(int32(len(open))))]
	w.Spawn(loc.Row, loc.Col)
}

// Returns error if it couldn't add monster
func (w *Room) AddMonster(loc Location, m monster.Type) error {
	if res := w.MonsterCanOccupy(loc); res != Move {
		return fmt.Errorf("Monster can't occupy %v: %v", loc, res)
	}
	w.monsters[loc] = append(w.monsters[loc], m)
	return nil
}

func (w *Room) SpawnMonster(m monster.Type) (Location, error) {
	var open []Location
	for row := 0; row < w.Rows(); row++ {
		for col := 0; col < w.Cols(); col++ {
			loc := Loc(row, col)
			if w.MonsterCanOccupy(loc) == Move {
				open = append(open, loc)
			}
		}
	}
	if len(open) == 0 {
		panic("Cannot find any place to spawn monster")
	}
	loc := open[int(rand.Int31n(int32(len(open))))]
	return loc, w.AddMonster(loc, m)
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
	if !w.InBounds(loc) {
		return OutOfBounds
	}
	if !w.TileAt(loc).Passable() {
		return Impassable
	}
	return Move
}

func (w *Room) MonsterCanOccupy(loc Location) MovementResult {
	if !w.InBounds(loc) {
		return OutOfBounds
	}
	if !w.TileAt(loc).Passable() {
		return Impassable
	}
	// Make sure monster doesn't spawn next to a door. That's really evil
	for doorLoc := range w.doors {
		if doorLoc.ManhattanDist(loc) <= 1 {
			return NextToDoor
		}
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
	case NextToDoor:
		return "Next to a door"
	default:
		panic(fmt.Sprintf("unknown movement result type %v", m))

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

func RandomWalledRoom(doorMap map[DoorDir]bool) *Room {
	rows := MinRows + int(rand.Int31n(MaxRows-MinRows))
	cols := MinCols + int(rand.Int31n(MaxCols-MinCols))
	log.Printf("Creating room of dimensions %d rows %d cols", rows, cols)
	w := WalledRoom(rows, cols)
	for dir := range doorMap {
		loc := w.LocForDoor(dir)
		w.SetTile(loc, LockedDoor)
	}
	return w
}

func RandomRoom(doorMap map[DoorDir]bool) *Room {
	propIsland := float32(0.3)
	if rand.Float32() < propIsland {
		return IslandRoom(doorMap)
	}
	return RandomWalledRoom(doorMap)
}

func (w *Room) LocForDoor(d DoorDir) Location {
	switch d {
	case North:
		return Loc(0, w.Cols()/2)
	case South:
		return Loc(w.Rows()-1, w.Cols()/2)
	case West:
		return Loc(w.Rows()/2, 0)
	case East:
		return Loc(w.Rows()/2, w.Cols()-1)
	}
	panic(fmt.Sprintf("unknown door location %v", d))
}

func (w *Room) HasDoorTile(d DoorDir) bool {
	tile := w.TileAt(w.LocForDoor(d))
	return tile == UnlockedDoor || tile == LockedDoor
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// An Island in the middle of the room surrounded by water
func IslandRoom(doors map[DoorDir]bool) *Room {
	r := RandomWalledRoom(doors)
	// Flood room with water
	for row := 1; row < r.Rows()-1; row++ {
		for col := 1; col < r.Cols()-1; col++ {
			r.SetTile(Loc(row, col), Water)
		}
	}

	// Build an island in middle
	middleRow := r.Rows() / 2
	middleCol := r.Cols() / 2
	for row := middleRow - 1; row <= middleRow+1; row++ {
		for col := middleCol; col <= middleCol+1; col++ {
			r.SetTile(Loc(row, col), Floor)
		}
	}
	// Build bridges from doors to middle
	if len(doors) == 0 {
		panic("no doors given")
	}
	for dir := range doors {
		loc := r.LocForDoor(dir)
		log.Printf("placing bridge to connect %v with island, at location %v\n", dir, loc)
		// Build from the middle up to the location
		minRow := min(loc.Row, middleRow)
		maxRow := max(loc.Row, middleRow)
		minCol := min(loc.Col, middleCol)
		maxCol := max(loc.Col, middleCol)
		log.Printf("min row %d max row %d min col %d max col %d\n", minRow, maxRow, minCol, maxCol)
		for row := minRow; row <= maxRow; row++ {
			for col := minCol; col <= maxCol; col++ {
				bridgeLoc := Loc(row, col)
				log.Printf("bridge loc? %v\n", bridgeLoc)

				if r.TileAt(bridgeLoc) == Water {
					log.Printf("placing bridge at %d %d", row, col)
					r.SetTile(bridgeLoc, Bridge)
				}
			}
		}
	}
	return r
}
