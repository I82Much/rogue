package rogue

type Player struct {
	MaxLife int
	Life    int
}

func NewPlayer(life int) *Player {
	return &Player{
		MaxLife: life,
		Life:    life,
	}
}

func (p *Player) IsDead() bool {
	return p.Life <= 0
}

func (p *Player) Damage(life int) {
	p.Life -= life
}
