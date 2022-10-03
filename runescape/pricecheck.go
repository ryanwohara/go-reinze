package runescape

import irc "github.com/thoj/go-ircevent"

func PriceCheck(irccon *irc.Connection) {
	irccon.Privmsg("#rshelp", "+ge ^(flax|battlestaff)$|zulrah|ar ar|((ne|al) spi)")
}
