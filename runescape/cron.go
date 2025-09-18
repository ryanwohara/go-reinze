package runescape

import (
	"database/sql"
	"fmt"
	"time"

	irc "github.com/thoj/go-ircevent"
)

func RunscapeCronHandler(irccon *irc.Connection, database *sql.DB) {
	if time.Now().Minute() == 0 {
		if time.Now().Hour()%4 == 0 {
			PriceCheck(irccon)
		}

		CheckNews(irccon, database)
	}

	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(database)
}
