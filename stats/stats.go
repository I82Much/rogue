package stats

import (
	"fmt"
)

// The extra that's published when player dies or wins
type Stats struct {
	LettersTyped int
	Hits int
	CompletedWords int
	MonstersDefeated int
	Rounds int
}

// Returns accuracy string (hits out of attempts, and %)
func (s *Stats) Accuracy() string {
	if s.LettersTyped == 0 {
		return "NA"
	}
	return fmt.Sprintf("%d / %d (%.2f%%)", s.Hits, s.LettersTyped, 100.0*float32(s.Hits)/float32(s.LettersTyped))
}

func (s *Stats) Add(s2 Stats) {
	s.LettersTyped += s2.LettersTyped
	s.Hits += s2.Hits
	s.CompletedWords += s2.CompletedWords
	s.Rounds += s2.Rounds
	s.MonstersDefeated += s2.MonstersDefeated
}

func (s *Stats) String() string {
	return fmt.Sprintf("Monsters: %d\nRounds: %d\nPhrases: %d\nAccuracy %s", s.MonstersDefeated, s.Rounds, s.CompletedWords, s.Accuracy())
}