package main

import (
	"database/sql"

	"reinze/news"

	irc "github.com/thoj/go-ircevent"
)

func cronHandler(irccon *irc.Connection, database *sql.DB, queue chan string) {
	news.CheckNews(database, irccon, queue)
}
