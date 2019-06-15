package runescape

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// GetUsersOnline returns a string
// including the number of users
// actively on OldSchool RuneScape.
func GetUsersOnline() string {
	return getUsersOnline()
}

func getUsersOnline() string {
	resp, err := http.Get("https://oldschool.runescape.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile("<p class='player-count'>There are currently ([0-9,]+) people playing!</p>")
	result := re.FindString(string(body))
	players := strings.Fields(result)[4]
	return players
}

func Matches(message string) string {
	response := ""
	if message == "players" {
		response = "There are currently " + GetUsersOnline() + " players online."
	}

	return response
}
