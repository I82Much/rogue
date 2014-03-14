package combat

import (
	"math/rand"
	"time"
)

var (
	// FIXME this is a hack
	playerWords = map[int][]string{
		0: []string{"foo", "baaaaaar", "cazaarararararar"},
		1: []string{"avast", "ye matey", "shiver me timbers"},
	}
)

func GetWords(round int) []AttackWord {
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
