package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/ryanwohara/reinze/runescape"
	irc "github.com/thoj/go-ircevent"
)

func main() {
	channels := os.Getenv("IRC_CHANNELS")
	hostname := os.Getenv("IRC_HOST") + ":" + os.Getenv("IRC_PORT")

	ircnick1 := os.Getenv("IRC_NICK")
	irccon := irc.IRC(ircnick1, ircnick1)
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: os.Getenv("IRC_HOST")}
	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Privmsgf("NickServ", "ID %s", os.Getenv("IRC_PASS"))
	})
	irccon.AddCallback("366", func(e *irc.Event) {})
	irccon.AddCallback("396", func(e *irc.Event) { irccon.Join(channels) })

	exported(irccon)

	go heartBeat(irccon)

	err := irccon.Connect(hostname)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	irccon.Loop()
}

type binFunc func(irccon *irc.Connection)

func exported(irccon *irc.Connection) {
	available := []binFunc{addInvite, addJoin, addNotice, addPrivmsg}
	for a := 0; a < len(available); a++ {
		handle(available[a], irccon)
	}
}

func handle(function binFunc, irccon *irc.Connection) {
	function(irccon)
}

func heartBeat(irccon *irc.Connection) {
	handleHeartBeat(irccon)

	for range time.Tick(time.Second * 60) {
		handleHeartBeat(irccon)
	}
}

func handleHeartBeat(irccon *irc.Connection) {
	fmt.Println(time.Now(), "Heartbeat")

	go runescape.RunscapeCronHandler(irccon)
	go cronHandler(irccon)
}
