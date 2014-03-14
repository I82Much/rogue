package combat

// This monster represents someone who speaks in l33t speak
// Lots of this comes from en.m.wikipedia.org/wiki/Leet
import (
	"math/rand"
	"strings"
)

var (
	haxorWords = []string{
		"l33t",
		"1337",
		"n00b",
		"haxor",
		"h4x0r",
		"pwnzor",
		"teh",
		"skillz",
		"w00t",
		"warez",
		"phreaking",
		"uber",
		"suxxor",
		"suxorz",
		"pwn3d",
		"Pr0n",
		"OMFG",
		"roxxorz",
	}
)

// randS randomly chooses from given strings
func randS(s ...string) string {
	return s[rand.Int31n(int32(len(s)))]
}

// If the word ends in given suffix, replaces it with one of the replacement strings.
func maybeReplaceSuffix(word, suffix string, replacements []string) string {
	if strings.HasSuffix(word, suffix) {
		lastIndex := strings.LastIndex(word, suffix)
		return word[0:lastIndex] + randS(replacements...)
	}
	return word
}

func makeHaxor(w string) string {
	s := func(s1 ...string) []string {
		return s1
	}
	replacementDict := map[string][]string{
		"e": s("3"),
		"t": s("7"),
		"l": s("1"),
		"o": s("0"),
	}
	// ed -> d, t, or 3d
	w = maybeReplaceSuffix(w, "ed", s("d", "t", "3d"))

	// and, anned, ant -> &
	w = maybeReplaceSuffix(w, "and", s("&"))
	w = maybeReplaceSuffix(w, "anned", s("&"))
	w = maybeReplaceSuffix(w, "ant", s("&"))

	// Make random replacements; sometimes the random choice is to leave the letter alone.
	for orig, replacementStrings := range replacementDict {
		replacement := randS(replacementStrings...)
		w = strings.Replace(w, orig, replacement, -1)
	}

	// Maybe add -zorz, age, or xor
	w = maybeReplaceSuffix(w, "", s("", "zorz", "age", "xor"))

	return w
}

func haxorWordFunc(round int) []string {
	// Just for fun we'll use the other monster's words too
	// FIXME if we have time convert the words from blogger etc
	return chooseNRandomly(haxorWords, round+1)
}
