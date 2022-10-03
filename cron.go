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

	go reddit.CheckPosts(database)

	go twitter.CheckPosts(database, irccon)

	go news.CheckNews(database, irccon)

	if time.Now().Minute() == 0 {
		go runescape.CheckNews(database, irccon)
	}

	defer database.Close()
}
