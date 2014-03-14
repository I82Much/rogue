package combat

import (
	"fmt"
	"log"
	"math/rand"
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

// phrases - the phrases to convert into the attacks
// words per minute - target number of words per minute. Assumes 5 characters per word,
// including spaces and punctuation (see wikipedia en.m.wikipedia.org/wiki/Words_per_minute)
// delay - how much extra time should be inserted between launches of attacks?
// TODO(ndunn):
func AttackWords(phrases []string, wordsPerMinute int, delay time.Duration) []*AttackWord {
	if wordsPerMinute <= 0 {
		panic(fmt.Sprintf("need wpm > 0; got %d", wordsPerMinute))
	}

	var totalDelay = delay
	var attacks []*AttackWord
	// characters per minute
	cpm := wordsPerMinute * 5
	for _, phrase := range phrases {
		chars := len(phrase)
		minutes := float32(chars) / float32(cpm)
		seconds := minutes / 60.0
		timeOnScreen := time.Duration(int(1000*seconds)) * time.Millisecond

		attack := NewWord(phrase, timeOnScreen, totalDelay)

		log.Printf("adding attack word %v", attack)

		attacks = append(attacks, &attack)

		// FIXME should be able to have some overlap on the screen
		totalDelay = time.Duration(totalDelay.Nanoseconds() + delay.Nanoseconds() + timeOnScreen.Nanoseconds())
	}
	return attacks
}

func chooseNRandomly(candidates []string, n int) []string {
	if n > len(candidates) {
		panic("not set up that way")
	}
	indices := rand.Perm(len(candidates))
	var res []string
	for i := 0; i < n; i++ {
		index := indices[i]
		res = append(res, candidates[index])
	}
	return res
}

func (w *AttackWord) Damage() int {
	return len(w.word)
}

func (w *AttackWord) IsVisible() bool {
	return time.Now().After(w.onScreen)
}
