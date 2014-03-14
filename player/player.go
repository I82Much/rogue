package player

type Stats struct {
	MonstersDefeated int
}

type Player struct {
	Name        string
	CurrentLife int
	MaxLife     int
	Level       int
	Stats       Stats
	MaxWPM      int
}

func WithName(n string) *Player {
	return &Player{
		Name:        n,
		CurrentLife: 100,
		MaxLife:     100,
		Level:       1,
		// FIXME
		MaxWPM: 60,
	}
}

func (p *Player) IsDead() bool {
	return p.CurrentLife <= 0
}

func (p *Player) Damage(life int) {
	p.CurrentLife -= life
}
