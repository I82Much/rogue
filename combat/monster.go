package combat

// TODO(ndunn): Figure out how to avoid duplication between monster and player.
type Monster struct {
	MaxLife int
	Life    int
	Words   []*AttackWord
}



func NewMonster(life int) *Monster {
	return &Monster{
		MaxLife: life,
		Life:    life,
	}
}

func (m *Monster) IsDead() bool {
	return m.Life <= 0
}

func (m *Monster) Damage(life int) {
	m.Life -= life
}
