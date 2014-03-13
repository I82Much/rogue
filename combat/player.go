package combat

import (
	"math/rand"
	"time"
)

type Player struct {
	MaxLife int
	Life    int
}

var (
	// FIXME this is a hack
	playerWords = map[int][]string{
		0: []string{"foo", "baaaaaar", "cazaarararararar"},
		1: []string{"avast", "ye matey", "shiver me timbers"},
	}
)

func NewPlayer(cur, max int) *Player {
	return &Player{
		MaxLife: max,
		Life:    cur,
	}
}

func (p *Player) IsDead() bool {
	return p.Life <= 0
}

func (p *Player) Damage(life int) {
	p.Life -= life
}

func (p *Player) GetWords(round int) []AttackWord {
	w := playerWords[round]
	var res []AttackWord
	for i, word := range w {
		attack := NewWord(word, time.Duration(rand.Int31n(10))*time.Second, time.Duration(0))
		col := columns[i%len(columns)]
		attack.Col = col
		res = append(res, attack)
	}
	return res
}
