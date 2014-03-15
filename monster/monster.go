package monster

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

func (t Type) Description() string{
	switch t {
	case Haxor:
		return "pwning haxor"
	case Scammer:
		return "gullible seeking scammer"
	case Spammer:
		return "desperate spammer"
	case Blogger:
		return "self aggrandizing blogger"
	case HaxorScammer:
		return "scamming hax0r"
	case HaxorSpammer:
		return "spamming hax0r"
	case HaxorBlogger:
		return "bl0gging hax0r"
	}
	return ""
}