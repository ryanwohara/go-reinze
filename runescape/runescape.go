package runescape

import (
	"fmt"
	"golang.org/x/text/language"
    "golang.org/x/text/message"
)

// GetUsersOnline returns a string
// including the number of users
// actively on OldSchool RuneScape.
func GetUsersOnline() []string {
	players := getUsersOnline()

	fmt.Println(players)

	p := message.NewPrinter(language.English)

	return []string{p.Sprintf("There are currently %d OSRS players and %d RS3 players online. (Total: %d)", players[0], players[1], players[2])}
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
