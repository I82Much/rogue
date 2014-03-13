package rogue

import (
	"fmt"
	"time"

	"github.com/I82Much/rogue/combat"
	"github.com/I82Much/rogue/dungeon"
	"github.com/I82Much/rogue/title"
	"github.com/I82Much/rogue/gameover"
)

type Game struct {
	// TODO(ndunn): lots of game state here

	curModule Module
}

const (
	rows = 16
	cols = 32
)

// TODO(ndunn): this needs to be pulled out of the world
func makeCombatModule() Module {
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
		combat.NewWord("Foo", time.Duration(1)*time.Second),
	}

	module := combat.NewModule(player, []*combat.Monster{m1, m2})
	return module
}

// TODO(ndunn): randomize
func makeDungeon() Module {
	room1 := dungeon.WalledRoom(rows, cols)
	room1.Spawn(rows/2, cols/2)
	room1.SpawnMonster()
	room1.SpawnMonster()

	// 2nd room to east of room 1
	room2 := dungeon.WalledRoom(rows, cols)

	// 3rd room to south of room 1
	room3 := dungeon.WalledRoom(rows, cols)

	world := dungeon.NewWorld(2, 2)
	world.Set(dungeon.Loc(0, 0), room1)
	world.Set(dungeon.Loc(0, 1), room2)
	world.Set(dungeon.Loc(1, 0), room3)

	// Set up doors between the rooms
	d1_2 := &dungeon.Door{
		From: room1,
		To:   room2,
	}
	// Door to east
	room1.SetDoor(dungeon.Loc(rows/2, cols-1), d1_2)
	d2_1 := &dungeon.Door{
		From: room2,
		To:   room1,
		Same: d1_2,
	}
	d1_2.Same = d2_1
	room2.SetDoor(dungeon.Loc(rows/2, 0), d2_1)

	d1_3 := &dungeon.Door{
		From: room1,
		To:   room3,
	}
	// Door to south
	room1.SetDoor(dungeon.Loc(rows-1, cols/2), d1_3)
	d3_1 := &dungeon.Door{
		From: room3,
		To:   room1,
		Same: d1_3,
	}
	d1_3.Same = d3_1
	room3.SetDoor(dungeon.Loc(0, cols/2), d3_1)
	return dungeon.NewModule(world)
}

func NewGame() *Game {
	d := title.NewModule()
	g := &Game{
		curModule: d,
	}
	d.AddListener(g)
	return g
}

func (g *Game) Start() {
	g.curModule.Start()
}

func (g *Game) Stop() {
	g.curModule.Stop()
}

// Listen handles the state transitions between the different modules.
func (g *Game) Listen(e string) {
	switch e {
		// Title screen
	case title.Start, gameover.Restart:
		g.Stop()
		// TODO ndunn fix me should be makeDUngeon
		c := makeCombatModule()
		//c := makeDungeon()
		c.AddListener(g)
		g.curModule = c
		g.Start()
		// Dungeon
	case dungeon.EnterCombat:
		g.Stop()
		c := makeCombatModule()
		c.AddListener(g)
		g.curModule = c
		g.Start()
		
		// Combat
	case combat.PlayerDied:
		g.Stop()
		c := gameover.NewModule()
		c.AddListener(g)
		g.curModule = c
		g.Start()

	case combat.AllMonstersDied:
		g.Stop()
		// TODO get loot
		fmt.Printf("Game over - you win")

	default:
		fmt.Errorf("unknown event: %v\n", e)
		g.Stop()
	}
}
