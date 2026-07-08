package runescape

import (
	"database/sql"
	"sync"
	"time"

	irc "github.com/thoj/go-ircevent"
)

var (
	lastRunMutex sync.Mutex
	// Starts at boot time so the first run lands on the next hour
	// boundary instead of spamming the channel on every restart.
	lastRun = time.Now()
)

func RunscapeCronHandler(irccon *irc.Connection, database *sql.DB) {
	now := time.Now()

	lastRunMutex.Lock()
	run := shouldRunHourly(lastRun, now)
	if run {
		lastRun = now
	}
	lastRunMutex.Unlock()

	if !run {
		return
	}

	if now.Hour()%4 == 0 {
		PriceCheck(irccon)
	}

	CheckNews(irccon, database)
}

func shouldRunHourly(last time.Time, now time.Time) bool {
	return now.Truncate(time.Hour).After(last.Truncate(time.Hour))
}
