package combat

import (
	"time"
)

// The crux of the game involves typing words or phrases that the enemies are sending at you.

// TODO(ndunn): if I had more time I'd make clever attack phrases etc based on who the enemy is. e.g. if I'm facing
// off against scammer, the phrases would be combatting the scammer. But probably will just end up using same phrases

// The difficulty is based on how fast the words move off the screen; the faster they move, the less time the player has to react or
// complete typing the word.


type AttackWord struct {
	word    string
	spelled []rune
	// How much fraction of time has elapsed for this word? Will render differently; e.g. attack could be going up towards
	// the monsters, defense down towards player
	proportion   float64
	initialDelay time.Duration
	// When was it first rendered
	onScreen time.Time
	duration time.Duration
	Col      Column
}

func NewWord(word string, dur time.Duration, initialDelay time.Duration) AttackWord {
	return AttackWord{
		word:         word,
		proportion:   0.0,
		initialDelay: initialDelay,
		onScreen:     time.Now().Add(initialDelay),
		duration:     dur,
		Col:          Left,
	}
}

func (w *AttackWord) Damage() int {
	return len(w.word)
}

func (w *AttackWord) IsVisible() bool {
	return time.Now().After(w.onScreen)
}
