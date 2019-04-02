package main

import (
	"fmt"

	"crypto/tls"

	irc "github.com/thoj/go-ircevent"
)

const channel = "#go"
const serverssl = "irc.swiftirc.net:6697"

func main() {
	ircnick1 := "blatiblat"
	irccon := irc.IRC(ircnick1, "IRCTestSSL")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("366", func(e *irc.Event) {})
	export(irccon)
	err := irccon.Connect(serverssl)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	irccon.Loop()
}

func addTest(irccon *irc.Connection) {
	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		if event.Nick == "Dragon" {
			if event.Message() == "TEST" {
				irccon.Privmsgf(event.Arguments[0], "Test over SSL successful\n")
			}
		}
	})
}

func addInvite(irccon *irc.Connection) {
	irccon.AddCallback("INVITE", func(event *irc.Event) {
		if event.Nick == "Dragon" {
			irccon.Join(event.Arguments[1])
		}
	})
}

type binFunc func(irccon *irc.Connection)

func export(irccon *irc.Connection) {
	available := [2]binFunc{addInvite, addTest}
	for a := 0; a < len(available); a++ {
		available[a](irccon)
	}
}
