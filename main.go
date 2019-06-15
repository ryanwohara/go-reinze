package main

import (
	"crypto/tls"
	"fmt"
	irc "github.com/thoj/go-ircevent"
	"os"
)

const channel = "#reinze"
const serverssl = "irc.swiftirc.net:6697"

func main() {
	ircnick1 := "PiKick"
	irccon := irc.IRC(ircnick1, ircnick1)
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: false} // change to `true` if you really have to
	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Privmsgf("NickServ", "ID %s", os.Getenv("REINZE_PASS"))
	})
	irccon.AddCallback("366", func(e *irc.Event) {})
	irccon.AddCallback("396", func(e *irc.Event) { irccon.Join(channel) })
	exported(irccon)
	err := irccon.Connect(serverssl)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	irccon.Loop()
}

type binFunc func(irccon *irc.Connection)

func exported(irccon *irc.Connection) {
	available := []binFunc{addInvite, addPrivmsg, addNotice}
	for a := 0; a < len(available); a++ {
		handle(available[a], irccon)
	}
}

func handle(function binFunc, irccon *irc.Connection) {
	function(irccon)
}
