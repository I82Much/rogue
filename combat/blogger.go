package combat

var (
	// a mix of things I find myself doing with annoying things I see others doing. Some link-bait type stuff in the mix as well.
	bloggerPhrases = []string{
		"clearly",
		"absolutely",
		"in my opinion",
		"if only",
		"follow me",
		"repin",
		"reblog",
		"followers",
		"it's clear to me",
		"comments disabled",
		"fork this",
		"you won't believe",
		"find out what happens when",
		"everything changes",
		// These come from wikipedia:
		// en.m.wikipedia.org/wiki/Glossary_of_blogging
		"rss",
		"blawg",
		"guest post",
		"blogosphere",
		"spam",
		"spammer",
		"fisking",
		"live blogging",
		"moblog",
		"mommy blog",
		"permalink",
		"pingback",
		"vlog",
	}
)

func bloggerWordFunc(round int) []string {
	return chooseNRandomly(bloggerPhrases, round+1)
}
