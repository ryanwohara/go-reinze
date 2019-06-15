package runescape

// GetUsersOnline returns a string
// including the number of users
// actively on OldSchool RuneScape.
func GetUsersOnline() string {
	return "There are currently " + getUsersOnline() + " players online."
}

// Matches what `runescape`
// package commands will run
// based on given triggers.
func Matches(message string) string {
	response := ""
	if message == "players" {
		response = GetUsersOnline()
	}

	return response
}
