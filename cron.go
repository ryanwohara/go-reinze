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

	database.Ping()

	reddit.CheckPosts(database)

	twitter.CheckPosts(database, irccon)

	news.CheckNews(database, irccon)

	if time.Now().Minute() == 0 {
		runescape.CheckNews(database, irccon)
	}

	defer database.Close()

	return
}
