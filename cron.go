package main

import (
	"database/sql"

	"github.com/ryanwohara/reinze/news"
	irc "github.com/thoj/go-ircevent"
)

func cronHandler(irccon *irc.Connection) {
	database := Db()

	err := database.Ping()

	if err != nil {
		println("cron.go: " + err.Error())

		return
	}

	news.CheckNews(database, irccon)

	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			println("cron.go: " + err.Error())
		}
	}(database)
}
