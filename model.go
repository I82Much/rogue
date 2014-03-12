package rogue

import (
	"github.com/I82Much/rogue/combat"
	"github.com/I82Much/rogue/dungeon"
)

type state int32

const (
	dungeonPhase = iota
	combatPhase
	gameoverPhase
)

type Model struct {
	phase state

	dungeonModel *dungeon.Model
	combatModel  *combat.Model
}

func NewModel(dm *dungeon.Model, cm *combat.Model) *Model {
	return &Model{
		phase:        dungeonPhase,
		dungeonModel: dm,
		combatModel:  cm,
	}
}

// dungeon <-> combat -> gameover
/*
type Model struct {

}*/
