package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"os"
	"reinze/runescape"
	"time"

	irc "github.com/thoj/go-ircevent"
)

func main() {
	channels := os.Getenv("IRC_CHANNELS")
	hostname := os.Getenv("IRC_HOST") + ":" + os.Getenv("IRC_PORT")

	ircnick1 := os.Getenv("IRC_NICK")
	irccon := irc.IRC(ircnick1, ircnick1)
	if irccon == nil {
		return
	}
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
	database := Db()

	incoming := make(chan string, 100)

	go func() {
		for msg := range incoming {
			irccon.SendRawf("PRIVMSG %s", msg)
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		go handleHeartBeat(irccon, database, incoming)
		time.Sleep(60 * time.Second)
	}
}

func handleHeartBeat(irccon *irc.Connection, database *sql.DB, queue chan string) {
	fmt.Println(time.Now(), "Heartbeat")

	go runescape.RunscapeCronHandler(irccon, database)
	go cronHandler(database, queue)
}
