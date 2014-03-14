package combat

import (
	"fmt"
	"testing"
)

func TestMakeHaxor(t *testing.T) {
	// One run yields
	// hello -> h3110zorz
	// how are you -> h0w ar3 y0u
	// what's the weather like today -> wha7's 7h3 w3a7h3r 1ik3 70day
	// this game would be a lot better if I were better at making games -> 7his gam3 w0u1d b3 a 107 b3773r if I w3r3 b3773r a7 making gam3sage

	words := []string{
		"hello",
		"how are you",
		"what's the weather like today",
		"this game would be a lot better if I were better at making games",
		"band",
	}
	for _, word := range words {
		fmt.Printf("%v -> %v\n", word, makeHaxor(word))
	}

}
