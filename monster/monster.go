package monster

import (
	"fmt"
	"math/rand"
)

type Type string

const (
	Haxor   = "HAXOR"
	Scammer = "SCAMMER"
	Spammer = "SPAMMER"
	Blogger = "BLOGGER"
	// Hybrids
	HaxorScammer = "HAXOR_SCAMMER"
	HaxorSpammer = "HAXOR_SPAMMER"
	HaxorBlogger = "HAXOR_BLOGGER"
)

var (
	All = []Type {
		Haxor,
		Scammer,
		Spammer,
		Blogger,
		HaxorScammer,
		HaxorSpammer,
		HaxorBlogger,
	}
)

var (
	descriptions = map[Type][]string{
		Haxor: []string{
			"pwning hax0r",
			"1337 h4x0r(s) looking to pwn noobs like you",
			"script kiddies looking to wreck your computer",
		},
		Scammer: []string{
			"scammer(s) seeking someone to play the role of next of kin",
		},
		Blogger: []string{
			"self aggrandizing blogger(s)",
			"blogger(s) who want you to join their email newsletter",
			"blogger(s) who really would like you to follow them",
			"blogger(s) who have a really great kickstarter for you to back",
		},
		Spammer: []string{
			"desperate spammer(s) looking to sell you some viagra",
			"folk(s) wanting you to Call Now",
		},
		HaxorScammer: []string{
			"scamming hax0r(s)",
		},
		HaxorSpammer: []string{
			"spamming hax0r(s)",
		},
		HaxorBlogger: []string{
			"bl0gging hax0r(s)",
		},
	}
)

// randS randomly chooses from given strings
func randS(s []string) string {
	return s[rand.Int31n(int32(len(s)))]
}

func (t Type) Description(num int) string {
	return fmt.Sprintf("%d %s", num, randS(descriptions[t]))
}
