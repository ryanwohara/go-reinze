package main

import (
	irc "github.com/thoj/go-ircevent"
)

func addInvite(irccon *irc.Connection) {
	irccon.AddCallback("INVITE", func(event *irc.Event) {
		if event.Nick == "Dragon" {
			irccon.Join(event.Arguments[1])
		}
	})
}
