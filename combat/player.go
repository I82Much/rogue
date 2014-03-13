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
		0: []string{"foo", "baaaaaar", "bazaarararararar"},
		1: []string{"avast", "ye matey", "shiver me timbers"},
	}
)

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

func (p *Player) GetWords(round int) []AttackWord {
	w := playerWords[round]
	var res []AttackWord
	for _, word := range w {
		res = append(res, NewWord(word, time.Duration(rand.Int31n(10))*time.Second))
	}
	return res
}
