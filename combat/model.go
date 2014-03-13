package combat

import (
	"fmt"
	"time"

	"github.com/I82Much/rogue/event"
	"github.com/I82Much/rogue/math"
	termbox "github.com/nsf/termbox-go"
)

const (
	PlayerDied      = "PLAYER_DIED"
	AllMonstersDied = "MONSTERS_VANQUISHED"
)

type State string
const (
	// Player is attacking
	Attack State = "ATTACK"
	// Player is defending
	Defend State = "DEFEND"
)

type Model struct {
	Monsters []*Monster
	Player   *Player

	words     []*AttackWord
	listeners []event.Listener

	state State

	attempts       int
	hits           int
	completedWords int
	currentTyping  *AttackWord
}

func (m *Model) AddListener(d event.Listener) {
	m.listeners = append(m.listeners, d)
}

func (m *Model) Publish(e string) {
	for _, listener := range m.listeners {
		listener.Listen(e)
	}
}

func NewCombatModel(p *Player, m []*Monster) *Model {
	var allWords []*AttackWord
	for _, m1 := range m {
		allWords = append(allWords, m1.Words...)
	}
	return &Model{
		Monsters: m,
		Player:   p,
		words:    allWords,
		// Player starts off defending against an onslaught of attacks
		state: Defend,
	}
}

type AttackWord struct {
	word    string
	spelled []rune
	// How much fraction of time has elapsed for this word? Will render differently; e.g. attack could be going up towards
	// the monsters, defense down towards player
	proportion float64
	onScreen time.Time
	duration time.Duration
}

func (w *AttackWord) Damage() int {
	return len(w.word)
}

func NewWord(word string, dur time.Duration) *AttackWord {
	return &AttackWord{
		word:     word,
		proportion: 0.0,
		onScreen: time.Now(),
		duration: dur,
	}
}

func (c *Model) Words() []*AttackWord {
	return c.words
}

func (c *Model) CurrentlyTyping() *AttackWord {
	return c.currentTyping
}

// KillWord removes the word from model, meaning that's it vanquished
func (c *Model) KillWord(w *AttackWord) {
	// TODO(ndunn): score? update exp?
	for i, word := range c.words {
		if word == w {
			c.words = append(c.words[0:i], c.words[i+1:]...)
			return
		}
	}
	// The terminal gets really screwed up if we don't shut down termbox first
	termbox.Close()
	panic(fmt.Sprintf("couldn't find word %v", w))
}

func (c *Model) DamagePlayer(w *AttackWord) {
	c.Player.Damage(w.Damage())
}

func (c *Model) DamageMonster(w *AttackWord) {
	// Pick the first monster that's not dead
	for _, monster := range c.Monsters {
		if !monster.IsDead() {
			monster.Damage(w.Damage())
		}
	}
	
}

func (c *Model) State() State {
	return c.state
}

// Over determines if the fight is over. Meaning either all enemies are dead, or player is dead
func (c *Model) PublishEndEvents() {
	if c.Player.IsDead() {
		c.Publish(PlayerDied)
	}
	// If any monster is left, fight's not over
	for _, m := range c.Monsters {
		if !m.IsDead() {
			return
		}
	}
	c.Publish(AllMonstersDied)
}

func (c *Model) Update(typed []rune) {
	now := time.Now()
	for _, r := range typed {
		c.attempts++
		// Does this rune represent the first untyped letter of any of the candidates? If so it's a hit. If not it's a miss
		if c.currentTyping != nil {
			runes := []rune(c.currentTyping.word)
			if r == runes[len(c.currentTyping.spelled)] {
				c.hits++
				c.currentTyping.spelled = append(c.currentTyping.spelled, r)

				// Done the word
				if len(c.currentTyping.spelled) == len(c.currentTyping.word) {
					c.KillWord(c.currentTyping)
					c.completedWords++
					c.currentTyping = nil
				}
			}
		} else {
			// See if the rune matches first letter of one of our candidate words
			for _, word := range c.Words() {
				runes := []rune(word.word)
				if r == runes[len(word.spelled)] {
					c.hits++
					word.spelled = append(word.spelled, r)
					c.currentTyping = word
					break
				}
			}
		}
	}

	var toRemove []*AttackWord
	for _, word := range c.words {
		elapsed := now.Sub(word.onScreen)
		// What proportion (0..1.0) is complete
		word.proportion = math.DoMap(float64(elapsed.Nanoseconds()), 0.0, float64(word.duration.Nanoseconds()), 0, 1.0)

		// Inflict damage on the player
		if word.proportion >= 1.0 {
			if c.state == Attack {
				c.DamageMonster(word)
			} else if c.state == Defend {
				c.DamagePlayer(word)
			} else {
				panic("shouldn't reach here")
			}
			toRemove = append(toRemove, word)
		}
	}
	for _, word := range toRemove {
		c.KillWord(word)
	}

	c.PublishEndEvents()
}
