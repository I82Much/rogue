package main

import (
	combat "github.com/I82Much/rogue/combat"
	
	"time"
)

func main() {
	player := combat.NewPlayer(20)
	m1 := combat.NewMonster(10)
	m1.Words = []*combat.AttackWord{
		combat.NewWord("Hello", time.Duration(3)*time.Second),
		combat.NewWord("Supercalifragilistic", time.Duration(10)*time.Second),
	}
	m2 := combat.NewMonster(50)
	m2.Words = []*combat.AttackWord{
		combat.NewWord("World", time.Duration(1)*time.Second),
		combat.NewWord("Blah", time.Duration(2)*time.Second),
		combat.NewWord("BlahBlah", time.Duration(2)*time.Second),
		combat.NewWord("Aasdfasdfasdfasdfasdfasdfasdf", time.Duration(10)*time.Second),
		//		combat.NewWord("World", time.Duration(20)*time.Second),

		combat.NewWord("Foo", time.Duration(1)*time.Second),
	}

	model := combat.NewCombatModel(player, []*combat.Monster{m1, m2})
	view := &combat.CombatView{
		Model: model,
	}
	controller := &combat.CombatController{
		Model: model,
		View:  view,
	}
	controller.Run(time.Duration(33) * time.Millisecond)
}
