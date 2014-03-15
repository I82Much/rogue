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
	// Hybrids
	HaxorScammer = "HAXOR_SCAMMER"
	HaxorSpammer = "HAXOR_SPAMMER"
	HaxorBlogger = "HAXOR_BLOGGER"
)

var (
	wordGenMap = map[MonsterType]func(round int) []string{
		Haxor:        haxorWordFunc,
		Scammer:      scammerWordFunc,
		Spammer:      spammerWordFunc,
		Blogger:      bloggerWordFunc,
		HaxorScammer: haxorBloggerWordFunc,
		HaxorSpammer: haxorSpammerWordFunc,
		HaxorBlogger: haxorBloggerWordFunc,
	}

	// Typing the h4xor stuff is very difficult. Slow it way down. A value of .5 in this case means multiply
	// the base WPM by this amount.
	wpmAdjustment = map[MonsterType]float32{
		Haxor:        0.5,
		HaxorScammer: 0.7,
		HaxorSpammer: 0.7,
		HaxorBlogger: 0.7,
	}

	// The haxor words are short but hard to type. Give a bonus to damage to make up for it.
	// The value changes meaning depending on if we're in attack phase or defense phase. e.g. a value of 2.0
	// indicates that if player successfully types it in attack, it does 2x damage. If he fails to type it
	// in defense mode, he takes the reciprocal (1/2.0) = .5x as much damage.
	damageAdjustment = map[MonsterType]float32{
		Haxor:        2.0,
		HaxorScammer: 1.5,
		HaxorSpammer: 1.5,
		HaxorBlogger: 1.5,
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

// Adjust the WPM based on the monster type. Slow down the hard to type haxor for instance.
func adjustWpm(t MonsterType, baseWpm int) int {
	if proportion, ok := wpmAdjustment[t]; ok {
		return int(float32(baseWpm) * proportion)
	}
	return baseWpm
}

func damageProportion(t MonsterType) float32 {
	if prop, ok := damageAdjustment[t]; ok {
		return prop
	}
	return 1.0
}

func getWords(round int, t MonsterType, wpm int) []*AttackWord {
	words := wordGenMap[t](round)
	// TODO what should the delay be
	log.Printf("monster %v round %d words: %v", t, round, words)

	attackWords := AttackWords(words, adjustWpm(t, wpm), time.Duration(rand.Int31n(2000))*time.Millisecond)
	// Modify the damage of words
	for _, w := range attackWords {
		w.DifficultyBonus = damageProportion(t)
	}
	return attackWords
}

func (p *Monster) GetWords(round int) []*AttackWord {
	return getWords(round, p.Type, p.WordsPerMinute)
}
