package gameover

import (
	"github.com/I82Much/rogue/static"
	"github.com/I82Much/rogue/title"
)

const (

	// generated by 	http://www.network-science.de/ascii/
	// font: stop
	screen = `


    ______                       _____                   
   / _____)                     / ___ \                  
  | /  ___  ____ ____   ____   | |   | |_   _ ____  ____ 
  | | (___)/ _  |    \ / _  )  | |   | | | | / _  )/ ___)
  | \____/( ( | | | | ( (/ /   | |___| |\ V ( (/ /| |    
   \_____/ \_||_|_|_|_|\____)   \_____/  \_/ \____)_|    
                                                       
                                                

		  Press 'r' to restart or 'q' to quit
		  Or
		  'e' for easy
		  'm' for medium
		  'h' for hard
		  
`

	Quit    = "QUIT"
	Restart = "RESTART_GAME"
)

func NewModule() *static.Module {
	return static.NewModule(screen, map[rune]string{
		'r': Restart,
		'R': Restart,
		'q': Quit,
		'Q': Quit,
		'e': title.Easy,
		'E': title.Easy,
		'm': title.Medium,
		'M': title.Medium,
		'h': title.Hard,
		'H': title.Hard,
		'i': title.Insane,
		'I': title.Insane,
		's': title.Stenographer,
		'S': title.Stenographer,
	})
}
