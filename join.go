package main

import (
	"fmt"

	irc "github.com/thoj/go-ircevent"
)

func addJoin(irccon *irc.Connection) {
	irccon.AddCallback("JOIN", func(event *irc.Event) {
		fmt.Println("~~>", event.Arguments[0])
		if event.Arguments[0] == "#asdfghj" {
			irccon.Privmsg(event.Arguments[0], "foo")
		}
	})
}
