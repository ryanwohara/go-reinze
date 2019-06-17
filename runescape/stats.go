package runescape

import (
	"fmt"
	"strings"
	"regexp"
)

// Stats will check all the stats
// of a given RuneScape username.
func Stats(message string) string {
	return track(message)
}

func stats(message string) string {
	fmt.Println(message)
	return "no data for you"
}

func getStats(rsn string) string {
	url := "http://services.runescape.com/m=hiscore_oldschool/hiscorepersonal.ws?user1=" + strings.Replace(rsn, " ", "+", -1)
	return url
}

func getSkill(match string) string {
	match = regexp.Regsub("\d+$", ", match) // TODO: make this valid
	skill_regex := [("((tota|overal)l|stats|oa)", "Overall"),
					("at(t(ack)?|k)", "Attack"),
					("def(en[cs]e)?", "Defence"),
					("str(ength)?", "Strength"),
					("h(p|it(point)?s)", "Hitpoints"),
					("rang(ed?|ing)", "Ranged"),
					("pray(er)?", "Prayer"),
					("mag(e|ic)", "Magic"),
					("cook(ing)?", "Cooking"),
					("w(c|ood(cut(ting)?))", "Woodcutting"),
					("fletch(ing)?", "Fletching"),
					("fish(ing)?", "Fishing"),
					("f(m|iremak(e|ing))", "Firemaking"),
					("craft(ing)?", "Crafting"),
					("(smith|smelt)(ing)?", "Smithing"),
					("min(e|ing)", "Mining"),
					("herb(l(aw|ore))?", "Herblore"),
					("agil(ity)?", "Agility"),
					("thie(f|v(e|ing))", "Thieving"),
					("slay(er|ing)?", "Slayer"),
					("farm(er|ing)?|rak(e|in[\"g]?)", "Farming"),
					("r(c|unecraft)(ing)?", "Runecraft"),
					("hunt(er|ing)?", "Hunter"),
					("con(struct(ion|ing)?)?", "Construction")]
	for regex in skill_regex:
		if regexp.MatchString(regex[0], match): // TODO: is this valid??
			return regex[1]