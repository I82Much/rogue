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
	// The damage is multiplied by this factor to offset how hard it is to type some phrases (e.g. 1337)
	DifficultyBonus float32
}

func NewWord(word string, dur time.Duration, initialDelay time.Duration) AttackWord {
	return AttackWord{
		word:            word,
		proportion:      0.0,
		initialDelay:    initialDelay,
		onScreen:        time.Now().Add(initialDelay),
		duration:        dur,
		Col:             Left,
		DifficultyBonus: 1.0,
	}
}

// multiply(2.0, time.Duration) yeidls a duration twice as long
func multiply(factor float64, dur time.Duration) time.Duration {
	ns := dur.Nanoseconds()
	newNs := int(factor * float64(ns))
	return time.Duration(newNs) * time.Nanosecond
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

	log.Printf("attack words phrases %v wpm %d delay %v", phrases, wordsPerMinute, delay)

	var totalDelay = delay
	var attacks []*AttackWord
	// characters per minute
	cpm := wordsPerMinute * 5
	for _, phrase := range phrases {
		chars := len(phrase)
		minutes := float32(chars) / float32(cpm)
		seconds := minutes * 60.0
		// Don't change this to just seconds or you'll get truncation to zero. Not what we want.
		timeOnScreen := time.Duration(int(1000*seconds)) * time.Millisecond

		// This is a hack. Oh well it makes the game more fun
		if chars < 6 {
			timeOnScreen = multiply(2, timeOnScreen)
		} else if chars < 8 {
			timeOnScreen = multiply(1.5, timeOnScreen)
		}

		attack := NewWord(phrase, timeOnScreen, totalDelay)

		log.Printf("adding attack word %v; total time on screen %v", attack, timeOnScreen.Seconds())

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

func (w *AttackWord) DamageToPlayer() int {
	// e.g. high difficulty words get a bonus, which means that it does less damage to player
	return int(float32(w.Damage()) / w.DifficultyBonus)
}

func (w *AttackWord) DamageToMonster() int {
	return int(float32(w.Damage()) * w.DifficultyBonus)
}

func (w *AttackWord) Damage() int {
	return len(w.word)
}

func (w *AttackWord) IsVisible() bool {
	return time.Now().After(w.onScreen)
}
