package runescape

import (
	"regexp"
)

// GetUsersOnline returns a string
// including the number of users
// actively on OldSchool RuneScape.
func GetUsersOnline() string {
	return "There are currently " + getUsersOnline(getOsRsPlayerData) + " players online in OSRS."
}

// Matches what `runescape`
// package commands will run
// based on given triggers.
func Matches(message string) string {
	response := ""
	if message == "players" {
		response = GetUsersOnline()
	} else if len(regexp.MustCompile("stats?|o(ver)?a(ll)?").FindString(message)) > 0 {
		response = Stats(message)
	} else if message == "track" {
		response = Track(message)
	}

	return response
}
