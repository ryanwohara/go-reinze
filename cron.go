package main

import (
	"database/sql"
	"fmt"

	"github.com/ryanwohara/reinze/news"
	irc "github.com/thoj/go-ircevent"
)

func cronHandler(irccon *irc.Connection, database *sql.DB) {
	news.CheckNews(database, irccon)

	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(database)
}
