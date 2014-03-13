package combat

import (
	"math/rand"
	"time"
)

// TODO(ndunn): Figure out how to avoid duplication between monster and player.
type Monster struct {
	MaxLife int
	Life    int
	Words   []AttackWord
}

var (
	// FIXME this is a hack
	words = map[int][]string{
		0: []string{"grunt", "hRRUUUUNT", "jello"},
		1: []string{"no way", "jose", "as.."},
	}
)

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

func (p *Monster) GetWords(round int) []AttackWord {
	w := words[round]
	var res []AttackWord
	for _, word := range w {
		res = append(res, NewWord(word, time.Duration(rand.Int31n(10))*time.Second))
	}
	return res
}
