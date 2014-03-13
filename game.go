package rogue

import (
	"fmt"
	"log"
	//	"time"

	"github.com/I82Much/rogue/combat"
	"github.com/I82Much/rogue/dungeon"
	"github.com/I82Much/rogue/gameover"
	"github.com/I82Much/rogue/player"
	"github.com/I82Much/rogue/title"
)

type Game struct {
	curModule Module
	dungeonModule *dungeon.Controller
	player *player.Player
}

const (
	rows = 16
	cols = 32
)

// TODO(ndunn): this needs to be pulled out of the world
func makeCombatModule() Module {
	player := combat.NewPlayer(100)
	m1 := combat.NewMonster(5)
	/*m1.Words = []*combat.AttackWord{
		combat.NewWord("Hello", time.Duration(3)*time.Second),
		combat.NewWord("Supercalifragilistic", time.Duration(2)*time.Second),
	}*/
	m2 := combat.NewMonster(5)
	/*m2.Words = []*combat.AttackWord{
		combat.NewWord("World", time.Duration(1)*time.Second),
		combat.NewWord("Blah", time.Duration(2)*time.Second),
		combat.NewWord("BlahBlah", time.Duration(2)*time.Second),
		combat.NewWord("Aasdfasdfasdfasdfasdfasdfasdf", time.Duration(3)*time.Second),
		combat.NewWord("Foo", time.Duration(1)*time.Second),
	}*/

	module := combat.NewModule(player, []*combat.Monster{m1, m2})
	return module
}

// TODO(ndunn): randomize
func makeDungeon(p *player.Player) *dungeon.Controller {
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
	return dungeon.NewModule(world, p)
}

func NewGame() *Game {
	d := title.NewModule()
	g := &Game{
		curModule: d,
		player: &player.Player{
			Name: "Player 1",
		},
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
	log.Printf("got event %v", e)
	switch e {
	// Title screen
	case title.Start, gameover.Restart:
		g.Stop()
		dm := makeDungeon(g.player)
		dm.AddListener(g)
		g.dungeonModule = dm
		g.curModule = dm
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
		g.curModule = g.dungeonModule
		// Ugh.
		g.dungeonModule.ReplaceMonsterWithPlayer()
		g.Start()
		// TODO get loot

	default:
		fmt.Errorf("unknown event: %v\n", e)
		g.Stop()
	}
}
