package rogue

import (
	
	"fmt"
	"time"
	
	termbox "github.com/nsf/termbox-go"
)

type CombatModel struct {
	Monsters []*Monster
	Player   *Player

	words []*AttackWord

	attempts       int
	hits           int
	completedWords int
	currentTyping  *AttackWord
}

func NewCombatModel(p *Player, m []*Monster) *CombatModel {
	var allWords []*AttackWord
	for _, m1 := range m {
		allWords = append(allWords, m1.Words...)
	}
	return &CombatModel{
		Monsters: m,
		Player:   p,
		words:    allWords,
	}
}

type AttackWord struct {
	word    string
	spelled []rune
	// TODO(ndunn): The row seems like part of view.. not sure though
	row      int
	maxRows  int
	onScreen time.Time
	duration time.Duration
}

func (w *AttackWord) Damage() int {
	return len(w.word)
}

func NewWord(word string, dur time.Duration) *AttackWord {
	return &AttackWord{
		word:     word,
		maxRows:  25,
		onScreen: time.Now(),
		duration: dur,
	}
}

func (c *CombatModel) Words() []*AttackWord {
	return c.words
}

func (c *CombatModel) CurrentlyTyping() *AttackWord {
	return c.currentTyping
}

// Over determines if the fight is over. Meaning either all enemies are dead, or player is dead
func (c *CombatModel) Over() bool {
	if c.Player.IsDead() {
		return true
	}
	// If any monster is left, fight's not over
	for _, m := range c.Monsters {
		if !m.IsDead() {
			return false
		}
	}
	return true
}

// KillWord removes the word from model, meaning that's it vanquished
func (c *CombatModel) KillWord(w *AttackWord) {
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

func (c *CombatModel) DamagePlayer(w *AttackWord) {
	c.Player.Damage(w.Damage())
}

func (c *CombatModel) Update(typed []rune) {
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

	for _, word := range c.words {
		elapsed := now.Sub(word.onScreen)
		row := int(doMap(elapsed.Seconds(), 0.0, word.duration.Seconds(), 0, float64(word.maxRows-1)))
		word.row = row
		
		// Inflict damage on the player
		if row > word.maxRows {
			c.DamagePlayer(word)
			c.KillWord(word)
		}
	}
}
