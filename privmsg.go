package main

import (
	"reinze/runescape"

	irc "github.com/thoj/go-ircevent"
)

func addPrivmsg(irccon *irc.Connection) {
	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		handlePrivmsg(irccon, event)
	})
}

func handlePrivmsg(irccon *irc.Connection, event *irc.Event) {
	if event.Nick != "Dragon" {
		return
	}

	message := event.Message()

	if len(message) < 2 {
		return
	}

	input := message[0]
	msg := message[1:]

	output := runescape.Matches(msg)

	if len(output) == 0 {
		return
	}

	for i := 0; i < len(output); i++ {
		if string(input) == "-" {
			irccon.Notice(event.Nick, output[i])
		} else if string(input) == "+" {
			irccon.Privmsg(event.Arguments[0], output[i])
		}
	}
}
