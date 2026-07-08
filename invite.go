package main

import (
	irc "github.com/thoj/go-ircevent"
)

func addInvite(irccon *irc.Connection) {
	irccon.AddCallback("INVITE", func(event *irc.Event) {
		handleInvite(irccon, event)
	})
}

func handleInvite(irccon *irc.Connection, event *irc.Event) {
	if event.Nick == "Dragon" && len(event.Arguments) > 1 {
		irccon.Join(event.Arguments[1])
	}
}
