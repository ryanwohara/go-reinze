package runescape

// GetUsersOnline returns a string
// including the number of users
// actively on OldSchool RuneScape.
func GetUsersOnline() string {
	return getUsersOnline()
}

// Matches what `runescape`
// package commands will run
// based on given triggers.
func Matches(message string) string {
	response := ""
	if message == "players" {
		response = "There are currently " + GetUsersOnline() + " players online."
	}

	return response
}
