package combat

import (
	"log"
	"math/rand"
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

func getWords(round int, t MonsterType, wpm int) []*AttackWord {
	words := wordGenMap[t](round)
	// TODO what should the delay be
	log.Printf("monster %v round %d words: %v", t, round, words)
	
	// The l33t speak is realllly hard to type. Give the player a break
	if t == Haxor {
		wpm = int(0.7 * float32(wpm))
	}
	return AttackWords(words, wpm, time.Duration(rand.Int31n(2000)) * time.Millisecond)
}

func (p *Monster) GetWords(round int) []*AttackWord {
	return getWords(round, p.Type, p.WordsPerMinute)
}
