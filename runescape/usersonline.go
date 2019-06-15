package runescape

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

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
	re := regexp.MustCompile("<p class='player-count'>There are currently ([0-9,]+) people playing OSRS!</p>")
	result := re.FindString(string(body))
	players := strings.Fields(result)[4]
	return players
}
