package main

import (
	"context"
	"database/sql"

	"github.com/ryanwohara/reinze/news"
	irc "github.com/thoj/go-ircevent"
)

func cronHandler(irccon *irc.Connection, database *sql.DB) {
	ctx, cancel := context.WithCancel(context.Background())
	if err := database.PingContext(ctx); err != nil {
		println("cron.go: " + err.Error())

		return
	}

	news.CheckNews(database, irccon)

	defer cancel()
}
