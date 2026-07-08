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

	err := irccon.Connect(hostname)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	// SendRaw writes to a channel that only exists once Connect
	// has succeeded, so the heartbeat must not start before it.
	go heartBeat(irccon)

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

	incoming := make(chan string, 1000)

	go func() {
		for msg := range incoming {
			safely("send queue", func() { irccon.SendRawf("%s", msg) })
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		for {
			handleHeartBeat(irccon, database, incoming)
			time.Sleep(90 * time.Second)
		}
	}()

	for {
		go safely("runescape cron", func() { runescape.RunscapeCronHandler(irccon, database) })
		time.Sleep(120 * time.Second)
	}
}

func handleHeartBeat(irccon *irc.Connection, database *sql.DB, queue chan string) {
	fmt.Println(time.Now(), "Heartbeat")

	go safely("news cron", func() { cronHandler(database, queue) })
}

// A panic in any of the cron goroutines would otherwise
// take down the entire bot.
func safely(name string, function func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(time.Now(), "recovered panic in", name, ":", r)
		}
	}()

	function()
}
