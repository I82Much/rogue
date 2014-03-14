package combat

import (
	"log"
	"time"
)

// TODO(ndunn): Figure out how to avoid duplication between monster and player.
type Monster struct {
	MaxLife        int
	Life           int
	WordsPerMinute int
	Words          []AttackWord
	Type           MonsterType
}

type MonsterType string

const (
	Haxor   = "HAXOR"
	Scammer = "SCAMMER"
	Spammer = "SPAMMER"
	Blogger = "BLOGGER"
)

var (
	wordGenMap = map[MonsterType]func(round int) []string{
		Haxor:   haxorWordFunc,
		Scammer: scammerWordFunc,
		Spammer: spammerWordFunc,
		Blogger: bloggerWordFunc,
	}
)

func NewMonster(life int, wordsPerMinute int, t MonsterType) *Monster {
	return &Monster{
		MaxLife:        life,
		Life:           life,
		WordsPerMinute: wordsPerMinute,
		Type:           t,
	}
}

func (m *Monster) IsDead() bool {
	return m.Life <= 0
}

func (m *Monster) Damage(life int) {
	m.Life -= life
}

func (p *Monster) GetWords(round int) []*AttackWord {
	words := wordGenMap[p.Type](round)
	// TODO what should the delay be
	log.Printf("monster %v round %d words: %v", *p, round, words)
	return AttackWords(words, p.WordsPerMinute, time.Duration(0))
}
