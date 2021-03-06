package combat

import (
	"fmt"
	"log"
	"time"

	"github.com/I82Much/rogue/event"
	"github.com/I82Much/rogue/math"
	"github.com/I82Much/rogue/player"
	"github.com/I82Much/rogue/stats"
	termbox "github.com/nsf/termbox-go"
)

const (
	PlayerDied      = "PLAYER_DIED"
	AllMonstersDied = "MONSTERS_VANQUISHED"
)

type State string

type Column string

const (
	EnemyDescription State = "DESCRIBE_ENEMY"
	EnteringAttack   State = "ENTERING_ATTACK"
	// Player is attacking
	Attack          State = "ATTACK"
	EnteringDefense       = "ENTERING_DEFENSE"
	// Player is defending
	Defense State = "DEFENSE"

	Left   Column = "LEFT"
	Right  Column = "RIGHT"
	Center Column = "CENTER"
)

var (
	columns = []Column{Left, Center, Right}

	// TODO(ndunn): this could shorten each time
	interRoundTime         = time.Duration(750) * time.Millisecond
	initialDescriptionTime = time.Duration(1500) * time.Millisecond
)

type Model struct {
	Monsters []*Monster
	Player   *player.Player

	words     []*AttackWord
	listeners []event.Listener

	state State

	attempts         int
	hits             int
	completedWords   int
	monstersDefeated int
	currentTyping    *AttackWord

	// which round of combat
	round int

	timeOfTransition time.Time
}

func (m *Model) MakeStats() stats.Stats {
	return stats.Stats{
		LettersTyped:     m.attempts,
		Hits:             m.hits,
		CompletedWords:   m.completedWords,
		Rounds:           m.round,
		MonstersDefeated: m.monstersDefeated,
	}
}

func (m *Model) AddListener(d event.Listener) {
	m.listeners = append(m.listeners, d)
}

func (m *Model) Publish(e string, extras interface{}) {
	for _, listener := range m.listeners {
		listener.Listen(e, extras)
	}
}

func NewCombatModel(p *player.Player, m []*Monster) *Model {
	return &Model{
		Monsters:         m,
		Player:           p,
		state:            EnemyDescription,
		timeOfTransition: time.Now().Add(initialDescriptionTime),
	}
}

func (c *Model) Words() []*AttackWord {
	return c.words
}

func (c *Model) CurrentlyTyping() *AttackWord {
	return c.currentTyping
}

func (c *Model) getPlayerAttackWords(round int) []*AttackWord {

	// Alternate between all of the alive monsters. e.g. round 1 pick first. Round 2, second. etc.
	var aliveMonsters []*Monster
	for _, m := range c.Monsters {
		if !m.IsDead() {
			aliveMonsters = append(aliveMonsters, m)
		}
	}
	if len(aliveMonsters) == 0 {
		panic("Invariant broken- there should be at least one alive monster")
	}
	return aliveMonsters[round%len(aliveMonsters)].GetWords(round)
}

func (c *Model) getAttackWords() []*AttackWord {
	var allWords []*AttackWord
	if c.state == Attack {
		allWords = c.getPlayerAttackWords(c.round)
	} else if c.state == Defense {
		for i, m := range c.Monsters {
			if m.IsDead() {
				continue
			}
			for _, word := range m.GetWords(c.round) {
				word := word
				// Change the column that the word falls from based on which monster it is.
				word.Col = columns[i%len(columns)]
				allWords = append(allWords, word)
			}
		}
	}
	return allWords
}

// KillWord removes the word from model, meaning that's it vanquished
func (c *Model) KillWord(w *AttackWord) {
	// TODO(ndunn): score? update exp?
	if c.currentTyping == w {
		log.Printf("no longer typing %v", c.currentTyping)
		c.currentTyping = nil
	}
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
	c.Player.Damage(w.DamageToPlayer())
}

func (c *Model) DamageMonster(w *AttackWord) {
	// Pick the first monster that's not dead
	for _, monster := range c.Monsters {
		if !monster.IsDead() {
			monster.Damage(w.DamageToMonster())
			// did this attack kill?
			if monster.IsDead() {
				c.monstersDefeated++
			}
			return
		}
	}
}

func (c *Model) State() State {
	return c.state
}

// Over determines if the fight is over. Meaning either all enemies are dead, or player is dead
func (c *Model) PublishEndEvents() {
	if c.Player.IsDead() {
		c.Publish(PlayerDied, c.MakeStats())
	}
	// If any monster is left, fight's not over
	for _, m := range c.Monsters {
		if !m.IsDead() {
			return
		}
	}
	c.Publish(AllMonstersDied, c.MakeStats())
}

// maybeTransition potentially shifts the model into another phase. e.g. after all the words are done in combat round,
// we enter the EnteringDefense round.
// The transitions are from
// EnteringDefense -> Defense -> EnteringAttack -> Attack -> EnteringDefense and on and on.
func (c *Model) maybeTransition() {
	if c.state == EnemyDescription && time.Now().After(c.timeOfTransition) {
		c.state = EnteringAttack
		c.timeOfTransition = time.Now().Add(interRoundTime)
	} else if c.state == Defense && len(c.words) == 0 {
		log.Println("defense -> entering attack")
		c.state = EnteringAttack
		c.timeOfTransition = time.Now().Add(interRoundTime)
	} else if c.state == Attack && len(c.words) == 0 {
		log.Println("attack -> entering defense")
		c.state = EnteringDefense
		c.timeOfTransition = time.Now().Add(interRoundTime)
		c.round++
	} else if c.state == EnteringDefense && time.Now().After(c.timeOfTransition) {
		log.Println("entering defense -> defense")
		c.state = Defense
		if len(c.words) != 0 {
			panic(fmt.Sprintf("invariant violated: should have had 0 words; had %v", c.words))
		}
		c.words = c.getAttackWords()
		//fmt.Printf("entered defense with words %v", c.words)
	} else if c.state == EnteringAttack && time.Now().After(c.timeOfTransition) {
		log.Println("entering attack -> attack")
		c.state = Attack
		if len(c.words) != 0 {
			panic(fmt.Sprintf("invariant violated: should have had 0 words; had %v", c.words))
		}
		c.words = c.getAttackWords()
		//fmt.Printf("entered attack with words %v", c.words)
	}
}

func (c *Model) Update(typed []rune) {
	// FIXME ndunn take this out
	//c.Publish(AllMonstersDied)

	now := time.Now()
	for _, r := range typed {
		c.attempts++
		// Does this rune represent the first untyped letter of any of the candidates? If so it's a hit. If not it's a miss
		if c.currentTyping != nil {
			log.Printf("Currently typing: %v", *c.currentTyping)
			runes := []rune(c.currentTyping.word)
			if r == runes[len(c.currentTyping.spelled)] {
				c.hits++
				c.currentTyping.spelled = append(c.currentTyping.spelled, r)
				log.Printf("spelled %v", c.currentTyping.spelled)

				// Done the word
				if len(c.currentTyping.spelled) == len(c.currentTyping.word) {
					log.Printf("finished %v", c.currentTyping.word)
					if c.state == Attack {
						// Finished typing the word - inflict damage if in attack mode
						c.DamageMonster(c.currentTyping)
					}
					c.KillWord(c.currentTyping)
					c.completedWords++
				}
			}
		} else {
			log.Printf("new letter %v", r)
			// See if the rune matches first letter of one of our candidate words
			for _, word := range c.Words() {
				runes := []rune(word.word)
				if r == runes[len(word.spelled)] {
					log.Printf("started typing %v", word.word)
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

		if word.proportion >= 1.0 {
			// Inflict damage on the player if in defense mode
			if c.state == Defense {
				c.DamagePlayer(word)
			}
			toRemove = append(toRemove, word)
		}
	}
	for _, word := range toRemove {
		c.KillWord(word)
	}

	// Transition phases
	c.maybeTransition()
	c.PublishEndEvents()
}
