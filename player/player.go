package player

type Stats struct {
	MonstersDefeated int
}

type Life int

type Player struct {
	Name string
	Current Life
	Max Life
	Level int
	Stats Stats
}