package player

type Stats struct {
	MonstersDefeated int
}

type Life int

type Player struct {
	Name    string
	Current Life
	Max     Life
	Level   int
	Stats   Stats
}

func WithName(n string) *Player {
	return &Player{
		Name:    n,
		Current: 100,
		Max:     100,
		Level:   1,
	}
}
