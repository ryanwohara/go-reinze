package main

import (
	"fmt"

	"crypto/tls"

	"github.com/ryanwohara/reinze/runescape"
	irc "github.com/thoj/go-ircevent"
)

const channel = "#reinze"
const serverssl = "irc.swiftirc.net:6697"

func main() {
	ircnick1 := "PiKick"
	irccon := irc.IRC(ircnick1, "IRCTestSSL")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("366", func(e *irc.Event) {})
	exported(irccon)
	err := irccon.Connect(serverssl)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	irccon.Loop()
}

func addInvite(irccon *irc.Connection) {
	irccon.AddCallback("INVITE", func(event *irc.Event) {
		if event.Nick == "Dragon" {
			irccon.Join(event.Arguments[1])
		}
	})
}

func addPrivmsg(irccon *irc.Connection) {
	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		if event.Nick == "Dragon" {
			msg := event.Message()[1:]
			input := event.Message()[0]

			output := runescape.Matches(msg)

			if len(output) == 0 {
				return
			}
			if string(input) == "-" {
				irccon.Notice(event.Nick, output)
			} else if string(input) == "+" {
				irccon.Privmsgf(event.Arguments[0], output)
			}
		}
	})
}

type binFunc func(irccon *irc.Connection)

func exported(irccon *irc.Connection) {
	available := []binFunc{addInvite, addPrivmsg}
	for a := 0; a < len(available); a++ {
		handle(available[a], irccon)
	}
}

func handle(function binFunc, irccon *irc.Connection) {
	function(irccon)
}
