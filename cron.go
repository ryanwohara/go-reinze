package main

import (
	"time"

	"github.com/ryanwohara/reinze/news"
	"github.com/ryanwohara/reinze/runescape"
	"github.com/ryanwohara/reinze/social/reddit"
	"github.com/ryanwohara/reinze/social/twitter"
	irc "github.com/thoj/go-ircevent"
)

func cronHandler(irccon *irc.Connection) {
	database := db()

	err := database.Ping()
	if err != nil {
		println("cron.go: " + err.Error())

		return
	}

	reddit.CheckPosts(database)

	twitter.CheckPosts(database, irccon)

	news.CheckNews(database, irccon)

	if time.Now().Minute() == 0 {
		runescape.CheckNews(database, irccon)

		if time.Now().Hour()%4 == 0 {
			runescape.PriceCheck(irccon)
		}
	}

	defer database.Close()
}
