package runescape

import irc "github.com/thoj/go-ircevent"

func PriceCheck(irccon *irc.Connection) {
	irccon.Privmsg("#rshelp", "+money -l 1k")
}
