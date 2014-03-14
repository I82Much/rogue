package combat

var (
	spammerPhrases = []string{
		"viagra",
		"enlargement",
		"nice tax credit",
		"work from home",
		"one weird trick",
		"great information",
		// Some of these come from TODO insert link
		// Hubspot.com's "The Ultimate List of Email SPAM Trigger Words", by Karen Rubin
		// posted January 11, 2012
		"meet singles",
		"order",
		"buy direct",
		"double your",
		"earn extra cash",
		"cents on the dollar",
		"mortage",
		"fees",
		"unsecured credit",
		"cards accepted",
		"full refund",
		"short pick",
		"avoid bankruptcy",
		"dormant",
		"stop",
		"deer friend",
		"click here",
		"marketing",
		"subscribe",
		"this isn't spam",
		"we hate spam",
		"lose weight",
		"xanax",
		"pharmacy",
		"100% free",
		"call",
		"no risk",
		"no obligation",
		"best rates",
		"compare",
		"get",
		"fax",
		"free",
		"free quote",
		"free sample",
		"certified",
		"congratulations",
		"satisfaction guaranteed",
		"act now",
		"offer expires",
		"limited time",
		"legal",
		"stainless steel",
	}
)

func spammerWordFunc(round int) []string {
	return chooseNRandomly(spammerPhrases, round+1)
}
