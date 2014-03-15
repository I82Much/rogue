package combat

import (
	"log"
	"math/rand"
	"time"

	"github.com/I82Much/rogue/monster"
)

// TODO(ndunn): Figure out how to avoid duplication between monster and player.
type Monster struct {
	MaxLife        int
	Life           int
	WordsPerMinute int
	Words          []AttackWord
	Type           monster.Type
}

var (
	wordGenMap = map[monster.Type]func(round int) []string{
		monster.Haxor:        haxorWordFunc,
		monster.Scammer:      scammerWordFunc,
		monster.Spammer:      spammerWordFunc,
		monster.Blogger:      bloggerWordFunc,
		monster.HaxorScammer: haxorBloggerWordFunc,
		monster.HaxorSpammer: haxorSpammerWordFunc,
		monster.HaxorBlogger: haxorBloggerWordFunc,
	}

	// Typing the h4xor stuff is very difficult. Slow it way down. A value of .5 in this case means multiply
	// the base WPM by this amount.
	wpmAdjustment = map[monster.Type]float32{
		monster.Haxor:        0.5,
		monster.HaxorScammer: 0.3,
		monster.HaxorSpammer: 0.3,
		monster.HaxorBlogger: 0.3,
	}

	// The haxor words are short but hard to type. Give a bonus to damage to make up for it.
	// The value changes meaning depending on if we're in attack phase or defense phase. e.g. a value of 2.0
	// indicates that if player successfully types it in attack, it does 2x damage. If he fails to type it
	// in defense mode, he takes the reciprocal (1/2.0) = .5x as much damage.
	damageAdjustment = map[monster.Type]float32{
		monster.Haxor:        2.0,
		monster.HaxorScammer: 3.0,
		monster.HaxorSpammer: 3.0,
		monster.HaxorBlogger: 3.0,
	}
)

func NewMonster(life int, wordsPerMinute int, t monster.Type) *Monster {
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
func adjustWpm(t monster.Type, baseWpm int) int {
	if proportion, ok := wpmAdjustment[t]; ok {
		return int(float32(baseWpm) * proportion)
	}
	return baseWpm
}

func damageProportion(t monster.Type) float32 {
	if prop, ok := damageAdjustment[t]; ok {
		return prop
	}
	return 1.0
}

func getWords(round int, t monster.Type, wpm int) []*AttackWord {
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
