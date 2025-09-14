package runescape

import (
	"database/sql"
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	irc "github.com/thoj/go-ircevent"
)

// GetUsersOnline returns a string
// including the number of users
// actively on OldSchool RuneScape.
func GetUsersOnline() []string {
	players := getUsersOnline()

	fmt.Println(players)

	p := message.NewPrinter(language.English)

	message := p.Sprintf("There are currently %d OSRS players (%.2f%%%%) and %d RS3 players (%.2f%%%%) online. (Total: %d) (Total Registered Accounts: %d)", players[0], float64(players[0])/float64(players[2])*100, players[1], float64(players[1])/float64(players[2])*100, players[2], players[3])

	fmt.Println([]string{message})

	return []string{message}
}

// Matches what `runescape`
// package commands will run
// based on given triggers.
func Matches(message string) []string {
	response := []string{}
	if message == "players" {
		response = GetUsersOnline()
	}

	return response
}

// Check OSRS and RS3 news
// and update the topic of
// #rshelp appropriately.
func CheckNews(irccon *irc.Connection, database *sql.DB) {
	checkNews(irccon, database)
}
