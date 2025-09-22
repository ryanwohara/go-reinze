package main

import (
	"database/sql"

	"reinze/news"
)

func cronHandler(database *sql.DB, queue chan string) {
	news.CheckNews(database, queue)
}
