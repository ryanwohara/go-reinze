package main

import (
	"github.com/ryanwohara/reinze/runescape"
	"github.com/ryanwohara/reinze/social/reddit"
	irc "github.com/thoj/go-ircevent"
	"time"
)

func cronHandler(irccon *irc.Connection) {
	db := db()

	db.Ping()

	reddit.CheckPosts(db)

	if time.Now().Minute() == 0 {
		news := runescape.CheckNews()
		irccon.SendRawf("TOPIC #rshelp :%s", news)
	}

	return
}