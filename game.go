package rogue

import (
	"fmt"
	"log"
	//	"time"

	"github.com/I82Much/rogue/combat"
	"github.com/I82Much/rogue/dungeon"
	"github.com/I82Much/rogue/gameover"
	"github.com/I82Much/rogue/monster"
	"github.com/I82Much/rogue/player"
	"github.com/I82Much/rogue/stats"
	"github.com/I82Much/rogue/title"
)

type Game struct {
	curModule     Module
	dungeonModule *dungeon.Controller
	player        *player.Player
	playerWpm     int
}

const (
	EasyWpm         = 15
	MediumWpm       = 40
	HardWpm         = 70
	InsaneWpm       = 100
	StenographerWpm = 300
)

var (
	difficultyMap = map[string]int{
		title.Easy:         EasyWpm,
		title.Medium:       MediumWpm,
		title.Hard:         HardWpm,
		title.Insane:       InsaneWpm,
		title.Stenographer: StenographerWpm,
	}
)

func (g *Game) lifeForMonster(t monster.Type) int {
	// TODO use player's level
	// TODO change it based on the monster
	return 10
}

// TODO(ndunn): this needs to be pulled out of the world
func (g *Game) makeCombat(t []monster.Type) Module {
	if len(t) == 0 {
		panic("need >= 1 monster")
	}
	player := g.player
	var monsters []*combat.Monster
	for _, m := range t {
		life := g.lifeForMonster(m)
		m1 := combat.NewMonster(life, player.MaxWPM, m)
		monsters = append(monsters, m1)
	}
	module := combat.NewModule(player, monsters)
	return module
}

func makeDungeon(p *player.Player) *dungeon.Controller {
	return dungeon.NewModule(dungeon.RandomWorld(2, 2), p)
}

func NewGame() *Game {
	d := title.NewModule()
	wpm := MediumWpm
	g := &Game{
		curModule: d,
		playerWpm: wpm,
		player:    player.WithName("Player 1", wpm),
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

func (g *Game) setWpm(wpm int) {
	g.playerWpm = wpm
	g.player.MaxWPM = wpm
}

func (g *Game) updateStats(stats stats.Stats) {
	g.player.Stats.Add(stats)
}

func (g *Game) restart() {
	g.Stop()
	g.player = player.WithName("Player 1", g.playerWpm)

	/*
		cm := g.makeCombat([]monster.Type{monster.Scammer})
		cm.AddListener(g)
		g.curModule = cm
		g.Start()*/

	dm := makeDungeon(g.player)
	dm.AddListener(g)
	g.dungeonModule = dm
	g.curModule = dm
	g.Start()
}

// Listen handles the state transitions between the different modules.
func (g *Game) Listen(e string, extra interface{}) {
	log.Printf("got event %v", e)
	switch e {
	case gameover.Restart:
		g.restart()
	// Title screen
	case title.Easy, title.Medium, title.Hard, title.Insane, title.Stenographer:
		g.setWpm(difficultyMap[e])
		g.restart()
	case dungeon.EnterCombat:
		g.Stop()
		types := extra.([]monster.Type)
		c := g.makeCombat(types)
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
		g.updateStats(extra.(stats.Stats))
		g.Stop()
		// Check to see if we've completed the game
		g.dungeonModule.ReplaceMonsterWithPlayer()
		g.dungeonModule.MaybeUnlockCurrentRoom()

		if g.dungeonModule.HasWon() {
			win := gameover.NewWinModule(g.player)
			win.AddListener(g)
			g.curModule = win
		} else {
			g.curModule = g.dungeonModule
		}
		g.Start()

	default:
		fmt.Errorf("unknown event: %v\n", e)
		g.Stop()
	}
}
