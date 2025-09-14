package runescape

import (
	"database/sql"
	"time"

	irc "github.com/thoj/go-ircevent"
)

func RunscapeCronHandler(irccon *irc.Connection) {
	if time.Now().Minute() == 0 {
		if time.Now().Hour()%4 == 0 {
			PriceCheck(irccon)
		}

		CheckNews(irccon)
	}

	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			println("cron.go: " + err.Error())
		}
	}(database)
}
