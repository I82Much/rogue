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

func Render(w *game.World) {
	r := w.CurrentRoom()
	for row := 0; row < r.Rows(); row++ {
		for col := 0; col < r.Cols(); col++ {
			// col = x, row = y
			location := game.Loc(row, col)
			termbox.SetCell(col, row, r.RuneAt(location), termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
	termbox.HideCursor()
}

func main() {

	// Set up controller
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.HideCursor()

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

	Render(world)

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

	//

	// Main game loop
	for {

		// Read input
		event := termbox.PollEvent()
		switch event.Key {
		case termbox.KeyArrowUp:
			if res := world.MovePlayer(-1, 0); res == game.CreatureOccupying {
				controller.Run(time.Duration(33) * time.Millisecond)
			}
		case termbox.KeyArrowRight:
			world.MovePlayer(0, 1)
		case termbox.KeyArrowDown:
			world.MovePlayer(1, 0)
		case termbox.KeyArrowLeft:
			world.MovePlayer(0, -1)
			// Quit
		case termbox.KeyCtrlC:
			return
		}
		// Render world
		Render(world)
	}
}
