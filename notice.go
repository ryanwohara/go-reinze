package main

import (
	irc "github.com/thoj/go-ircevent"
)

func addNotice(irccon *irc.Connection) {
	irccon.AddCallback("NOTICE", func(event *irc.Event) {
		//
	})
}
