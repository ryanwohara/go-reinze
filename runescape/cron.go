package runescape

import (
	"database/sql"
	"time"

	irc "github.com/thoj/go-ircevent"
)

func RunscapeCronHandler(irccon *irc.Connection, database *sql.DB) {
	if time.Now().Minute() == 0 {
		if time.Now().Hour()%4 == 0 {
			PriceCheck(irccon)
		}

		CheckNews(database, irccon)
	}
}
