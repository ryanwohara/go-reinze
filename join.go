package main

import (
	"os"

	"reinze/greet"

	irc "github.com/thoj/go-ircevent"
)

func addJoin(irccon *irc.Connection) {
	irccon.AddCallback("JOIN", func(event *irc.Event) {
		handleJoin(irccon, event)
	})
}

func handleJoin(irccon *irc.Connection, event *irc.Event) {
	if len(event.Arguments) == 0 {
		return
	}

	if event.Arguments[0] == os.Getenv("GREET_CHANNEL") && os.Getenv("IRC_NICK") != event.Nick {
		greet.Greet(irccon, event.Arguments[0], event.Nick)
	}
}
