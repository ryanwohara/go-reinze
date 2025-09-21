package main

import (
	"os"

	"reinze/greet"

	irc "github.com/thoj/go-ircevent"
)

func addJoin(irccon *irc.Connection) {
	irccon.AddCallback("JOIN", func(event *irc.Event) {
		if event.Arguments[0] == os.Getenv("GREET_CHANNEL") && os.Getenv("IRC_NICK") != event.Nick {
			greet.Greet(irccon, event.Arguments[0], event.Nick)
			greet.Greet(irccon, event.Arguments[0], event.Nick)
		}
	})
}
