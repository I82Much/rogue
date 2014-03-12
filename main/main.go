package main

import (
	"time"

	combat "github.com/I82Much/rogue/combat"
	game "github.com/I82Much/rogue/dungeon"
	termbox "github.com/nsf/termbox-go"
)

const (
	rows = 16
	cols = 32
)

func makeDungeon() *game.Model {
	room1 := game.WalledRoom(rows, cols)
	room1.Spawn(rows/2, cols/2)
	room1.SpawnMonster()
	room1.SpawnMonster()

	// 2nd room to east of room 1
	room2 := game.WalledRoom(rows, cols)

	// 3rd room to south of room 1
	room3 := game.WalledRoom(rows, cols)

	world := game.NewWorld(2, 2)
	world.Set(game.Loc(0, 0), room1)
	world.Set(game.Loc(0, 1), room2)
	world.Set(game.Loc(1, 0), room3)

	// Set up doors between the rooms
	d1_2 := &game.Door{
		From: room1,
		To:   room2,
	}
	// Door to east
	room1.SetDoor(game.Loc(rows/2, cols-1), d1_2)
	d2_1 := &game.Door{
		From: room2,
		To:   room1,
		Same: d1_2,
	}
	d1_2.Same = d2_1
	room2.SetDoor(game.Loc(rows/2, 0), d2_1)

	d1_3 := &game.Door{
		From: room1,
		To:   room3,
	}
	// Door to south
	room1.SetDoor(game.Loc(rows-1, cols/2), d1_3)
	d3_1 := &game.Door{
		From: room3,
		To:   room1,
		Same: d1_3,
	}
	d1_3.Same = d3_1
	room3.SetDoor(game.Loc(0, cols/2), d3_1)

	return game.NewModel(world)
}

func makeCombatModel() *combat.Model {
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
	return model
}

func main() {
	// Set up controller
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.HideCursor()
	defer termbox.Close()

	dungeonModel := makeDungeon()
	dungeonView := game.NewView(dungeonModel)
	dungeonView.Render()
	dungeonController := game.NewController(dungeonModel, dungeonView)

	/*
		combatModel := makeCombatModel()
		combatView := &combat.View{
			Model: combatModel,
		}
		combatController := &combat.Controller{
			Model: combatModel,
			View:  combatView,
		}*/

	dungeonController.Run()
}
