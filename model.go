package rogue

type state int32

const (
	dungeon = iota
	combat
	gameover
)

// dungeon <-> combat -> gameover


